/*
Copyright © 2023 Alexander Orban <alexander.orban@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package _import

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"golang.org/x/sync/errgroup"
	"log"
	"math"
	"regexp"
	es "snowlastic-cli/pkg/es"
	orm "snowlastic-cli/pkg/orm"
	"sync"
	"time"
)

var segmenter string

// import/casesCmd represents the import/cases command
var casesSegmentedCmd = &cobra.Command{
	Use:   "cases-segmented",
	Short: "Index all Navex Cases contained in the snowflake database",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet("esClient") {
			return errors.New("elasticsearch was somehow not created by the `import` command prior to running `demos`")
		}
		return nil
	},
	RunE: runImport,
}

func init() {
	casesSegmentedCmd.Flags().StringVar(&segmenter, "by", "", "a field or SQL aggregating function used to split the import")
}

func runImport(cmd *cobra.Command, args []string) error {
	var (
		dbSchema  = "SQL_NAVEX"
		dbTable   = "" // there are multiple, which are handled in the GetQuery return
		indexName = "test"
		docType   = (&orm.Case{}).New()

		db    *sql.DB
		query = docType.GetQuery(dbSchema, dbTable)

		err error
	)

	db, err = generateDB(dbSchema)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var segments []interface{}
	if segmenter != "" {
		segments, err = getSegments(db, query, segmenter)
	} else {
		segments = []interface{}{"*"}
		err = nil
	}
	log.Printf("processing %d segments into %s\n", len(segments), indexName)
	if err != nil {
		return err
	}

	g := errgroup.Group{}
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))

	for _, seg := range segments {
		var segment interface{}
		segment = seg
		g.Go(func() error {
			var thisQuery string
			switch segment {
			case "*", "%", "all":
				thisQuery = query
			case nil:
				thisQuery = query + " WHERE " + segmenter + " IS NULL"
			default:
				thisQuery = query + " WHERE " + segmenter + "= " + quoteParam(segment)
			}

			var rowCount int64
			var numBatches float64
			rowCount, err = getRowCount(db, thisQuery)
			if err != nil {
				return err
			}
			numBatches = math.Ceil(float64(rowCount) / es.BulkInsertSize)

			bar := p.AddBar(int64(numBatches),
				//mpb.BarRemoveOnComplete(),
				mpb.PrependDecorators(
					//decor.Name(barName, decor.WC{W: len(barName) + 1, C: decor.DidentRight}),
					decor.Name(fmt.Sprintf("%v", segment), decor.WCSyncSpaceR),
					decor.CountersNoUnit("%5d/%5d", decor.WCSyncWidth),
				),
				mpb.AppendDecorators(
					decor.Percentage(decor.WC{W: 5}),
					decor.OnComplete(
						// ETA decorator with ewma age of 30
						decor.EwmaETA(decor.ET_STYLE_GO, 30, decor.WCSyncWidth), " done",
					),
				),
			)

			start := time.Now().UTC()
			var cases = make(chan orm.SnowlasticDocument, es.BulkInsertSize)
			g.Go(func() error {
				defer close(cases)
				rows, err := db.Query(thisQuery)
				if err != nil {
					return fmt.Errorf("problem with stmt.Query(%s): %s", segment, err)
				}

				for rows.Next() {
					var c orm.Case
					if err := c.ScanFrom(rows); err != nil {
						return fmt.Errorf("problem with scanning case %s: %s", c.CaseID, err)
					}
					cases <- &c
				}
				return nil
			})
			c, err := getElasticClient(ElasticsearchClientLocator)
			if err != nil {
				return fmt.Errorf("problem with getting elasticsearch clinet for segment %s: %s", segment, err)
			}
			batches := es.BatchEntities(cases, es.BulkInsertSize)
			numIndexed, numErrors, err := es.BulkImportWithMPB(c, batches, indexName, bar)
			if err != nil {
				return fmt.Errorf("problem with bulk importing for segment %s: %s", segment, err)
			}
			if numErrors > 0 {
				return errors.New(fmt.Sprintf(
					"%s:\tIndexed [%s] documents with [%s] errors in %s (%s docs/sec)",
					fmt.Sprintf("%s: %s=%v", indexName, segmenter, segment),
					humanize.Comma(int64(numIndexed)),
					humanize.Comma(int64(numErrors)),
					time.Since(start).Truncate(time.Millisecond),
					humanize.Comma(int64(1000.0/float64(time.Since(start)/time.Millisecond)*float64(numIndexed))),
				))
			}

			return nil
		})
	}
	return g.Wait()
}

func quoteParam(i interface{}) string {
	switch i.(type) {
	case string:
		return fmt.Sprintf("'%s'", i)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", i)
	case float32, float64:
		return fmt.Sprintf("%f", i)
	}
	return ""
}
func quoteField(i interface{}) string {
	switch i.(type) {
	case string:
		whitespace := regexp.MustCompile(`\s`).MatchString(i.(string))
		if whitespace {
			return fmt.Sprintf(`"%s"`, i)
		}
		return fmt.Sprintf("%s", i)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", i)
	case float32, float64:
		return fmt.Sprintf("%f", i)
	}
	return ""
}
