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
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"golang.org/x/sync/errgroup"
	"log"
	"math"
	"math/rand"
	"snowlastic-cli/pkg/es"
	orm "snowlastic-cli/pkg/orm"
	"sync"
	"time"
)

// purchaseOrdersCmd represents the import/purchaseOrders command
var purchaseOrdersSegmentedCmd = &cobra.Command{
	Use:   "purchaseOrders-segmented",
	Short: "Index all Purchase Orders contained in the snowflake database",
	Long:  ``,
	RunE:  runPOsSegmentedImport,
}

func init() {
	purchaseOrdersSegmentedCmd.Flags().StringVar(&segmenter, "by", "", "a field or SQL aggregating function used to split the import")
}

func runPOsSegmentedImport(cmd *cobra.Command, args []string) error {
	var (
		dbSchema  = "CMP"
		dbTable   = "PURCHASEORDERS_FLAGGED" // there are multiple, which are handled in the GetQuery return
		indexName = "purchaseorders"
		docType   = (&orm.PurchaseOrder{}).New()

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
		segments, err = getSegments(db, query, quoteField(segmenter))
	} else {
		segments = []interface{}{"*"}
		err = nil
	}
	log.Printf("processing %d segments into %s\n", len(segments), indexName)
	if err != nil {
		return err
	}

	g := errgroup.Group{}
	g.SetLimit(5)
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
				thisQuery = query + " WHERE " + quoteField(segmenter) + " IS NULL" + " LIMIT " + randBetween(3000, 6000)
			default:
				thisQuery = query + " WHERE " + quoteField(segmenter) + "= " + quoteParam(segment) + " LIMIT " + randBetween(3000, 6000)
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
					decor.CountersNoUnit("%5d/%5d ", decor.WCSyncWidth),
					decor.Percentage(decor.WC{W: 5}),
				),
				mpb.AppendDecorators(
					decor.OnComplete(
						// ETA decorator with ewma age of 30
						decor.EwmaETA(decor.ET_STYLE_GO, 30, decor.WCSyncWidth), " done",
					),
				),
			)

			start := time.Now().UTC()
			var cases = make(chan orm.SnowlasticDocument, es.BulkInsertSize)

			h := errgroup.Group{}
			h.Go(func() error {
				defer close(cases)
				rows, err := db.Query(thisQuery)
				if err != nil {
					return fmt.Errorf("problem with stmt.Query(%s): %s", segment, err)
				}

				for rows.Next() {
					var c = docType.New()
					if err := c.ScanFrom(rows); err != nil {
						return fmt.Errorf("problem with scanning case %s: %s", c.GetID(), err)
					}
					cases <- c
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

			return h.Wait()
		})
	}
	return g.Wait()
}

func randBetween(min, max int) string {
	r := rand.Intn(max-min) + min
	return fmt.Sprintf("%d", r)
}
