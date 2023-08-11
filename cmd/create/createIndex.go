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
package create

import (
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"snowlastic-cli/pkg/es"
)

var (
	isVendor   bool
	isCustomer bool
	isDemo     bool
	runAll     bool
	fromFile   string

	c *elasticsearch.Client
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !isVendor && !isCustomer && !isDemo && fromFile == "" {
			return errors.New("at least one flag is required by the index command")
		}
		if fromFile != "" && len(args) == 0 {
			return errors.New(`you must provide an index name when creating an anonymous index from a file
Usage: snowlastic-cli.exe create index --from ./path/to/settings.json <index name>

`)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			caCertPath string
			err        error
			caCert     []byte
			cfg        es.ElasticClientConfig
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

		// parse the given commands
		if runAll {
			fmt.Println("creating vendor index")
			if err != nil {
				return err
			}
			fmt.Println("creating customer index")
			if err != nil {
				return err
			}
			fmt.Println("creating demo index")
			err = indexDemo(c)
			if err != nil {
				return err
			}
		}
		if isVendor {
			fmt.Println("creating vendor index")
			if err != nil {
				return err
			}
		}
		if isCustomer {
			fmt.Println("creating customer index")
			if err != nil {
				return err
			}
		}
		if isDemo {
			fmt.Println("creating demo index")
			err = indexDemo(c)
			if err != nil {
				return err
			}
		}
		if fromFile != "" {
			fmt.Println("creating index from file")
			if err != nil {
				return err
			}
		}
		return err
	},
}

func init() {
	//createCmd.AddCommand(indexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//indexCmd.Flags().BoolVarP(&isVendor, "vendor", "v", false, "create a standard vendor index")
	//indexCmd.Flags().BoolVarP(&isCustomer, "customer", "c", false, "create a standard customer index")
	indexCmd.Flags().BoolVarP(&isDemo, "demo", "d", false, "create a standard demo index")
	indexCmd.Flags().BoolVarP(&runAll, "all", "", false, "create all standard indices")
	indexCmd.Flags().StringVarP(&fromFile, "from", "", "", "Create an anonymous index from a json file containing elasticsearch index settings")
	//indexCmd.Flags().BoolP("demo", "", false, "create a demo index")
}
