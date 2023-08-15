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
	"encoding/json"
	"errors"
	"fmt"
	icm_orm "github.com/alexander-orban/icm_goapi/orm"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"snowlastic-cli/pkg/es"
	"time"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			caCertPath string
			err        error
			caCert     []byte
			cfg        es.ElasticClientConfig
			c          *elasticsearch.Client

			filePath  string = args[0]
			indexName string = viper.GetString("identifier")
			As        []Anon

			docs       = make(chan icm_orm.ICMEntity, es.BulkInsertSize)
			numErrors  int64
			numIndexed int64
		)
		// generate the CA Certificate bytes needed for the elasticsearch Config
		caCertPath = viper.GetString("elasticCaCertPath")
		caCert, err = os.ReadFile(caCertPath)
		if err != nil {
			return err
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
		// Generate the client
		c, err = es.NewElasticClient(&cfg)
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
				docs <- icm_orm.ICMEntity(&As[i])
			}
			close(docs)
		}()

		batches := es.BatchEntities(docs, es.BulkInsertSize)
		numIndexed, numErrors, err = es.BulkImport(c, batches, indexName)
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
	fileCmd.MarkFlagRequired("identifier")
	_ = viper.BindPFlag("identifier", fileCmd.Flags().Lookup("identifier"))

	fileCmd.Flags().String("index", "", "index into which to import records")
	fileCmd.MarkFlagRequired("index")
	_ = viper.BindPFlag("index", fileCmd.Flags().Lookup("index"))
}

type m map[string]interface{}
type Anon struct {
	m
}

func (a *Anon) IsICMEntity() bool { return true }
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
