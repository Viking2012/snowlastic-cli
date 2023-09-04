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
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"golang.org/x/sync/errgroup"
	"log"
	"math"
	"math/rand"
	"os"
	"regexp"
	"snowlastic-cli/pkg/es"
	orm "snowlastic-cli/pkg/orm"
	"snowlastic-cli/pkg/snowflake"
	"sync"
	"time"
	"unicode"
)

var (
	ElasticsearchClientLocator = "esClient"

	segmenter     string
	givenSegments []string
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import documents into an elasticsearch index",
	Long: `A document is a representation of any kind of record. This tool allows
for importing data from pre-defined sources (such as snowflake tables/views) or 
from a json file containing a list of documents.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var (
			c   *elasticsearch.Client
			err error
		)
		c, err = generateElasticClient()
		if err != nil {
			return err
		}
		viper.Set(ElasticsearchClientLocator, c)

		if cmd.Flags().Lookup("in").Changed && !cmd.Flags().Lookup("by").Changed {
			return errors.New("you must specify a 'by' field in order to use the 'in' argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("import called")
	},
}

func Add() *cobra.Command {
	return importCmd
}

func init() {
	importCmd.AddCommand(demoCmd)
	importCmd.AddCommand(casesCmd)
	importCmd.AddCommand(casesSegmentedCmd)
	importCmd.AddCommand(fileCmd)
	importCmd.AddCommand(purchaseOrdersCmd)
	importCmd.AddCommand(purchaseOrdersSegmentedCmd)

	importCmd.PersistentFlags().StringVar(&segmenter, "by", "", "a field or SQL aggregating function used to split the import")
	importCmd.PersistentFlags().StringSliceVar(&givenSegments, "in", []string{}, "limit the import of the field defined in the 'by' argument to a comma seperated list of values")
}

// Elasticsearch utilities
func generateElasticClient() (*elasticsearch.Client, error) {
	var (
		err        error
		caCert     []byte
		caCertPath string
		cfg        es.ElasticClientConfig
		c          *elasticsearch.Client
	)

	// generate the CA Certificate bytes needed for the elasticsearch Config
	caCertPath = viper.GetString("elasticCaCertPath")
	caCert, err = os.ReadFile(caCertPath)
	if err != nil {
		return c, err
	}
	cfg = es.ElasticClientConfig{
		Addresses: []string{fmt.Sprintf(
			"https://%s:%s",
			viper.GetString("elasticUrl"),
			viper.GetString("elasticPort"),
		)},
		User:         viper.GetString("elasticUser"),
		Pass:         viper.GetString("elasticPassword"),
		ApiKey:       viper.GetString("elasticApiKey"),
		ServiceToken: viper.GetString("elasticServiceToken"),
		CaCert:       caCert,
	}
	return es.NewElasticClient(&cfg)
}
func getElasticClient(key string) (*elasticsearch.Client, error) {
	var (
		c  *elasticsearch.Client
		ok bool
	)
	c, ok = viper.Get(key).(*elasticsearch.Client)
	if !ok {
		return nil, errors.New("was not able to gather an elasticsearch client after being created by the `import` command")
	}
	return c, nil
}
func reportImport(prefix string, dur time.Duration, numIndexed, numErrors int64) error {
	if numErrors > 0 {
		return errors.New(fmt.Sprintf(
			"%s:\tIndexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			prefix,
			humanize.Comma(numIndexed),
			humanize.Comma(numErrors),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		))
	} else {
		log.Printf(
			"%s:\tSucessfully indexed [%s] documents in %s (%s docs/sec)",
			prefix,
			humanize.Comma(numIndexed),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	}
	return nil
}

// Database utilities
func generateDB(schema string) (*sql.DB, error) {
	return snowflake.NewDB(snowflake.Config{
		Account:   viper.GetString("snowflakeAccount"),
		Warehouse: viper.GetString("snowflakeWarehouse"),
		Database:  viper.GetString("snowflakeDatabase"),
		Schema:    schema,
		User:      viper.GetString("snowflakeUser"),
		Password:  viper.GetString("snowflakePassword"),
		Role:      viper.GetString("snowflakeRole"),
	})

}
func getRowCount(db *sql.DB, baseQuery string) (int64, error) {
	var rowCount int64

	countQuery := "SELECT COUNT(1) FROM (" + baseQuery + ")"
	rows, err := db.Query(countQuery)
	if err != nil {
		return rowCount, nil
	}
	rows.Next()
	err = rows.Scan(&rowCount)
	return rowCount, err
}
func getSegments(db *sql.DB, baseQuery string, segmenter string) ([]interface{}, error) {
	var segments []interface{}
	var segmentationQuery = "SELECT DISTINCT " +
		segmenter +
		" FROM (" +
		baseQuery +
		") ORDER BY " +
		segmenter
	rows, err := db.Query(segmentationQuery)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var segment interface{}
		err := rows.Scan(&segment)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}
	return segments, err
}

// importation utilities
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
		var needsQuotes bool
		var ret string
		re := regexp.MustCompile(`^(.*\()(.*)(\).*)`)
		match := re.FindStringSubmatch(`"Creation Date"`)
		if len(match) != 0 {
			needsQuotes, _ = needsQuoting(match[2])
			if needsQuotes {
				ret = fmt.Sprintf(`%s"%s"%s`, match[1], match[2], match[3])
			}
			ret = fmt.Sprintf(`%s%s%s`, match[1], match[2], match[3])
		} else {
			needsQuotes, _ = needsQuoting(i.(string))
			if needsQuotes {
				ret = fmt.Sprintf(`"%s"`, i)
			}
			ret = fmt.Sprintf("%s", i)
		}
		return ret
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", i)
	case float32, float64:
		return fmt.Sprintf("%f", i)
	}
	return ""
}
func needsQuoting(field string) (bool, error) {
	var (
		matched bool
		err     error
	)
	matched, err = regexp.Match(`^[A-Za-z_].*`, []byte(field))
	if err != nil {
		return true, err
	}
	if !matched {
		return true, nil
	}

	// contains any non-alphanumeric character
	matched, err = regexp.Match(".*[^A-Za-z0-9_].*", []byte(field))
	if err != nil {
		return true, err
	}
	if matched {
		return true, nil
	}

	// Is not all uppercase
	if !isUpper(field) {
		return true, nil
	}

	return false, nil
}
func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
func randBetween(min, max int) string {
	r := rand.Intn(max-min) + min
	return fmt.Sprintf("%d", r)
}

// Generalized Import Command
func runSegmentedImport(dbSchema, dbTable, indexName string, docType orm.SnowlasticDocument) error {
	var (
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
		if len(givenSegments) == 0 {
			segments, err = getSegments(db, query, quoteField(segmenter))
		} else {
			segments = make([]interface{}, len(givenSegments))
			for i := range givenSegments {
				segments[i] = givenSegments[i]
			}
		}
	} else {
		segments = []interface{}{"*"}
		err = nil
	}
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
				thisQuery = query + " WHERE " + quoteField(segmenter) + " = " + quoteParam(segment) + " LIMIT " + randBetween(3000, 6000)
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
			h.SetLimit(1)
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
					humanize.Comma(numIndexed),
					humanize.Comma(numErrors),
					time.Since(start).Truncate(time.Millisecond),
					humanize.Comma(int64(1000.0/float64(time.Since(start)/time.Millisecond)*float64(numIndexed))),
				))
			}

			return h.Wait()
		})
	}
	return g.Wait()
}
