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
	"log"
	"os"
	"snowlastic-cli/pkg/es"
)

type run struct {
	f    func(*elasticsearch.Client) error
	name string
}

var (
	isCase   bool
	isDemo   bool
	runAll   bool
	fromFile string

	c    *elasticsearch.Client
	runs []run
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Create an elasticsearch index",
	Long:  "",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().NFlag() == 0 {
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
			runs = []run{
				{f: indexDemo, name: "demo"},
				{f: indexCase, name: "case"},
			}
		}

		if isDemo && !runAll {
			runs = append(runs, run{f: indexDemo, name: "demo"})
		}
		if isCase && !runAll {
			runs = append(runs, run{f: indexCase, name: "case"})
		}
		if fromFile != "" && !runAll {
			fmt.Printf("creating index '%s' from file\n", args[0])
			err = indexFile(c, fromFile, args[0])
			if err != nil {
				return err
			}
		}

		for _, run := range runs {
			log.Printf("creating %s index", run.name)
			err := run.f(c)
			if err != nil {
				return nil
			}
		}

		return nil
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
	indexCmd.Flags().BoolVarP(&runAll, "all", "", false, "create all standard indices")

	indexCmd.Flags().BoolVarP(&isCase, "case", "c", false, "create a standard Navex case index")
	indexCmd.Flags().BoolVarP(&isDemo, "demo", "d", false, "create a standard demo index")

	indexCmd.Flags().StringVarP(&fromFile, "from", "", "", "Create an anonymous index from a json file containing elasticsearch index settings")
}
