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
	icm_encoding "github.com/alexander-orban/icm_goapi/encoding"
	icm_orm "github.com/alexander-orban/icm_goapi/orm"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"log"
	"math"
	"snowlastic-cli/pkg/es"
	"time"
)

// purchaseOrdersCmd represents the import/purchaseOrders command
var purchaseOrdersCmd = &cobra.Command{
	Use:   "purchaseOrders",
	Short: "Index all Purchase Orders contained in the snowflake database",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			dbQueryFrom string = "PURCHASEORDERS_FLAGGED"
			indexName          = "purchaseorders"
			docType            = icm_orm.PurchaseOrder{Flags: &icm_orm.PurchaseOrderFlags{}}

			db  *sql.DB
			esC *elasticsearch.Client

			docs = make(chan icm_orm.ICMEntity, es.BulkInsertSize)

			numErrors  int64
			numIndexed int64

			err error
		)

		db, err = generateDB("CMP")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var query = icm_encoding.MarshalToSelect(&docType, dbQueryFrom, false)
		var rowCount int64
		var countQuery = "SELECT COUNT(1) FROM (" + string(query) + ")"

		// Get row and batch counts
		rows, err := db.Query(countQuery)
		if err != nil {
			return err
		}
		rows.Next()
		err = rows.Scan(&rowCount)
		if err != nil {
			return err
		}
		numBatches := math.Ceil(float64(rowCount) / es.BulkInsertSize)

		start := time.Now().UTC()
		go func() {
			log.Println("reading POs from database")
			rows, err := db.Query(string(query))
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				var e icm_orm.ICMEntity
				if e, err = icm_orm.PurchaseOrderFromRow(rows); err != nil {
					log.Fatal(err)
				}
				docs <- e
			}
			close(docs)
		}()

		log.Println("connecting to elasticsearch")
		esC, err = getElasticClient(ElassticsearchClientLocator)
		if err != nil {
			return err
		}

		log.Println("batching cases")
		batches := es.BatchEntities(docs, es.BulkInsertSize)
		log.Println("indexing cases")
		numIndexed, numErrors, err = es.BulkImport(esC, batches, indexName, int64(numBatches))
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
				"Successfully indexed [%s] documents in %s (%s docs/sec)",
				humanize.Comma(int64(numIndexed)),
				dur.Truncate(time.Millisecond),
				humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
			)
		}
		return nil
	},
}

func init() {}
