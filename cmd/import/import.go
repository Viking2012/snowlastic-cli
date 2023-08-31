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
	"os"
	"snowlastic-cli/pkg/es"
	"snowlastic-cli/pkg/snowflake"
	"time"
)

var (
	ElasticsearchClientLocator string = "esClient"
	DatabaseConnectionLocator  string = "db"
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
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("import called")
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if viper.IsSet(DatabaseConnectionLocator) {
			db, err := getDb(DatabaseConnectionLocator)
			if err != nil {
				return err
			}
			err = db.Close()
			if err != nil {
				return err
			}
		}
		return nil
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
}

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

func generateDB(schema string) (*sql.DB, error) {
	log.Println("connecting to database")
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

func getDb(key string) (*sql.DB, error) {
	db, ok := viper.Get(key).(*sql.DB)
	if !ok {
		return nil, errors.New("was not able to gather a database connection after being created by the `import` command")
	}
	return db, nil
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
		return segments, err
	}
	for rows.Next() {
		var segment interface{}
		err := rows.Scan(&segment)
		if err != nil {
			return segments, err
		}
		segments = append(segments, segment)
	}
	return segments, err
}

func reportImport(dur time.Duration, numIndexed, numErrors int64) error {
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
			"Sucessfully indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	}
	return nil
}
