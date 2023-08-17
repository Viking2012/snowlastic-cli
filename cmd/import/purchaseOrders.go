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
	"github.com/spf13/viper"
	"log"
	"math"
	"os"
	"snowlastic-cli/pkg/es"
	"snowlastic-cli/pkg/snowflake"
	"time"
)

// purchaseOrdersCmd represents the import/purchaseOrders command
var purchaseOrdersCmd = &cobra.Command{
	Use:   "purchaseOrders",
	Short: "Index all Purchase Orders contained in the snowflake database",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err error

			caCert     []byte
			caCertPath string
			esCfg      es.ElasticClientConfig
			esC        *elasticsearch.Client

			db *sql.DB

			docs = make(chan icm_orm.ICMEntity, es.BulkInsertSize)

			indexName = "purchaseorders"

			numErrors  int64
			numIndexed int64
		)

		log.Println("connecting to database")
		var tmp = icm_orm.PurchaseOrder{Flags: &icm_orm.PurchaseOrderFlags{}}
		var query = icm_encoding.MarshalToSelect(&tmp, "PURCHASEORDERS_FLAGGED", false)
		db, err = snowflake.NewDB(snowflake.Config{
			Account:   viper.GetString("snowflakeAccount"),
			Warehouse: viper.GetString("snowflakeWarehouse"),
			Database:  viper.GetString("snowflakeDatabase"),
			Schema:    "CMP",
			User:      viper.GetString("snowflakeUser"),
			Password:  viper.GetString("snowflakePassword"),
			Role:      viper.GetString("snowflakeRole"),
		})
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var rowCount int64
		var countQuery = "SELECT COUNT(1) FROM (" + string(query) + ")"
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

		// generate the CA Certificate bytes needed for the elasticsearch Config
		log.Println("connecting to elasticsearch")
		caCertPath = viper.GetString("elasticCaCertPath")
		caCert, err = os.ReadFile(caCertPath)
		if err != nil {
			return err
		}
		esCfg = es.ElasticClientConfig{
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
		// Generate the client
		esC, err = es.NewElasticClient(&esCfg)
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

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// purchaseOrdersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// purchaseOrdersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
