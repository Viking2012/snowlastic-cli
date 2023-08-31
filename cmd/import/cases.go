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
	"log"
	"math"
	es "snowlastic-cli/pkg/es"
	orm "snowlastic-cli/pkg/orm"
	"time"
)

// import/casesCmd represents the import/cases command
var casesCmd = &cobra.Command{
	Use:   "cases",
	Short: "Index all Navex Cases contained in the snowflake database",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet("esClient") {
			return errors.New("elasticsearch was somehow not created by the `import` command prior to running `demos`")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			dbSchema  = "SQL_NAVEX"
			dbTable   = "" // there are multiple, which are handled in the GetQuery return
			indexName = "cases"
			docType   = orm.PurchaseOrder{}

			db *sql.DB
			c  *elasticsearch.Client

			docs = make(chan orm.SnowlasticDocument, es.BulkInsertSize)

			err error
		)

		db, err = generateDB(dbSchema)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var query = docType.GetQuery(dbSchema, dbTable)
		var rowCount int64
		var numBatches float64
		rowCount, err = getRowCount(db, query)
		if err != nil {
			return err
		}
		numBatches = math.Ceil(float64(rowCount) / es.BulkInsertSize)

		start := time.Now().UTC()
		go func() {
			log.Println("reading cases from database")
			rows, err := db.Query(query)
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				var c orm.Case
				if err := c.ScanFrom(rows); err != nil {
					log.Fatal(err)
				}
				docs <- &c
			}
			close(docs)
		}()

		// Get the generated elasticsearch client
		c, err = getElasticClient(ElasticsearchClientLocator)
		if err != nil {
			return err
		}
		batches := es.BatchEntities(docs, es.BulkInsertSize)
		log.Println("indexing cases")
		numIndexed, numErrors, err := es.BulkImport(c, batches, indexName, int64(numBatches))
		if err != nil {
			return err
		}

		dur := time.Since(start)
		if numErrors > 0 {
			return errors.New(fmt.Sprintf(
				"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
				humanize.Comma(int64(numIndexed)),
				humanize.Comma(int64(numErrors)),
				dur.Truncate(time.Millisecond),
				humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
			))
		} else {
			log.Printf(
				"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
				humanize.Comma(int64(numIndexed)),
				dur.Truncate(time.Millisecond),
				humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
			)
		}
		return nil
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// import/casesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// import/casesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
