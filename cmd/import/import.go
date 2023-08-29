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
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"snowlastic-cli/pkg/es"
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
		viper.Set("esClient", c)
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
	importCmd.AddCommand(fileCmd)
	importCmd.AddCommand(purchaseOrdersCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getElasticClient() (*elasticsearch.Client, error) {
	var (
		c  *elasticsearch.Client
		ok bool
	)
	c, ok = viper.Get("esClient").(*elasticsearch.Client)
	if !ok {
		return nil, errors.New("was not able to gather an elasticsearch client after being created by the `import` command")
	}
	return c, nil
}
