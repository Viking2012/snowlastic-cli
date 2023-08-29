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
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"snowlastic-cli/pkg/es"
)

var (
	settingsFlag string

	c *elasticsearch.Client
)

var defaultEntitySettings = map[string]string{
	"cases":          "./settings/esindex-cases.json",
	"purchaseorders": "./settings/esindex-purchaseorders.json",
	"demos":          "./settings/esindex-demos.json",
}

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index index-name [--settings ./path/to/settings.json]",
	Short: "Create an elasticsearch index",
	Long: `Create an elasticsearch index of the given name, optionally with explicit settings from a json file
Example: create index demos --settings ./settings/elastic-indices/demos.json

Indices created without the --settings flag will be created with the elasticsearch servers default settings.
The following ICM Entities have pre-defined settings which will override the server defaults, even without a --settings flag being provided:
  ICM Entity                          | Index name	   		| Example command
-----------------------------------------------------------------------------------------
- Navex cases 							cases				  create index cases
- Demonstrations and keyword testing 	demos				  create index demos
- Purchase Orders 						purchaseorders		  create index purchaseorders`,
	ValidArgs: []string{"cases", "demos", "purchaseorders"},
	Args:      cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			caCertPath string
			err        error
			caCert     []byte

			settingsFilepath string
		)
		// generate the CA Certificate bytes needed for the elasticsearch Config
		caCertPath = viper.GetString("elasticCaCertPath")
		caCert, err = os.ReadFile(caCertPath)
		if err != nil {
			return err
		}
		// Generate the client
		c, err = es.NewElasticClient(&es.ElasticClientConfig{
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
		})
		if err != nil {
			return err
		}

		settingsFilepath = settingsFlag
		for i := range cmd.ValidArgs {
			if args[0] == cmd.ValidArgs[i] {
				settingsFilepath = defaultEntitySettings[args[0]]
				break
			}
		}

		var b []byte
		if settingsFilepath != "" {
			fmt.Printf("creating index '%s' with settings file located at %s\n", args[0], settingsFilepath)
			b, err = os.ReadFile(settingsFilepath)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("creating index '%s' with default settings\n", args[0])
		}

		err = createIndex(c, b, args[0])
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	indexCmd.Flags().StringVarP(&settingsFlag, "settings", "", "", "Create an anonymous index from a json file containing elasticsearch index settings")
}

func createIndex(c *elasticsearch.Client, settings []byte, indexName string) error {
	res, err := c.Indices.Delete([]string{indexName})
	if err != nil {
		return fmt.Errorf("cannot delete index: %s", err)
	}
	if res.IsError() {
		log.Println("error when deleting index", res.String())
	} else {
		log.Println(res.String())
	}

	res, err = c.Indices.Create(indexName, c.Indices.Create.WithBody(bytes.NewReader(settings)))
	if err != nil {
		return fmt.Errorf("cannot create index: %s", err)
	}

	if res.IsError() {
		return fmt.Errorf("cannot create index, got an error response code: %s\n", res.String())
	}
	log.Printf("successfully created index %s\n", indexName)
	return nil
}
