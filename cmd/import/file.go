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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"math"
	"os"
	"snowlastic-cli/pkg/es"
	orm "snowlastic-cli/pkg/orm"
	"time"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Import a json file into an elasticsearch index",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet("esClient") {
			return errors.New("elasticsearch was somehow not created by the `import` command prior to running `demos`")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			indexName string = viper.GetString("identifier")

			filePath string = args[0]
			As       []Anon

			c *elasticsearch.Client

			docs       = make(chan orm.SnowlasticDocument, es.BulkInsertSize)
			numErrors  int64
			numIndexed int64

			err error
		)
		// generate the CA Certificate bytes needed for the elasticsearch Config
		c, err = getElasticClient(ElasticsearchClientLocator)
		if err != nil {
			return err
		}

		log.Println("reading file", filePath)
		b, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		log.Println("unmarshalling records")
		err = json.Unmarshal(b, &As)
		if err != nil {
			return err
		}

		start := time.Now().UTC()
		go func() {
			for i := range As {
				docs <- &As[i]
			}
			close(docs)
		}()

		batches := es.BatchEntities(docs, es.BulkInsertSize)
		numIndexed, numErrors, err = es.BulkImport(c, batches, indexName, int64(math.Ceil(float64(len(As))/es.BulkInsertSize)))
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
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	fileCmd.Flags().String("identifier", "id", "Fieldname for unique identifier of each record")
	_ = fileCmd.MarkFlagRequired("identifier")
	_ = viper.BindPFlag("identifier", fileCmd.Flags().Lookup("identifier"))

	fileCmd.Flags().String("index", "", "index into which to import records")
	_ = fileCmd.MarkFlagRequired("index")
	_ = viper.BindPFlag("index", fileCmd.Flags().Lookup("index"))
}

type m map[string]interface{}
type Anon struct {
	m
}

func (a *Anon) IsDocument() {}
func (a *Anon) GetID() string {
	var (
		idField string
		s       string
	)
	idField = viper.GetString("identifier")
	s = fmt.Sprintf("%s", a.m[idField])
	return s
}
func (a *Anon) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(a.m)
	return j, err
}
func (a *Anon) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	a.m = m
	return nil
}
func (a *Anon) GetQuery(string, string) string { return "" }
func (a *Anon) ScanFrom(rows *sql.Rows) error  { return nil }
func (a *Anon) New() orm.SnowlasticDocument    { return new(Anon) }
