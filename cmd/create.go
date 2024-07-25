/*
	package cmd

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
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	settings string
)

var indexConfigurations = make(map[string]map[string]string)
var knownIndices = make(map[string]string)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create index-name.... [--settings ./path/to/settings.json]",
	Short: "Create an elasticsearch index",
	Long: `Create an elasticsearch index of the given name, optionally with explicit settings from a json file
Example/Multiple Indices : create sap_checks sap_purchaseorders navex_cases demos 
Example/Explicit Settings: create demos --settings ./settings/esindex-demos.json

Indices created without the --settings flag will be created with the individual settings defined in the 
snowlastic-cli.yaml configuration file. If the --settings flag is provided with multiple index names, 
this file will be used to create each index.
Indices not defined in the snowlastic-cli.yaml configuration file and are created without a --settings flag will
be created without any explicit mapping (elasticsearch defaults to its own internal typing).`,
	//Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// initialize primary variables
		indexConfigurations = viper.Get("elasticIndices").(map[string]map[string]string)
		for indexName, indexPaths := range indexConfigurations {
			knownIndices[indexName] = indexPaths["path_to_index_settings"]
		}
		return
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			c   *elasticsearch.Client
			err error

			givenIndex       string
			settingsFilepath string
		)

		// Generate the client
		c, err = generateDefaultElasticClient()
		if err != nil {
			return err
		}
		if len(args) == 0 {
			for indexName, _ := range knownIndices {
				args = append(args, indexName)
			}
		}

		for i := range args {
			givenIndex = args[i]

			if isKnownIndex(knownIndices, givenIndex) && settings == "" {
				settingsFilepath = knownIndices[givenIndex]
			} else {
				settingsFilepath = settings
			}

			var b []byte
			if settingsFilepath != "" {
				log.Println("using settings defined in", settingsFilepath)
				b, err = os.ReadFile(settingsFilepath)
				if err != nil {
					return err
				}
			} else {
				log.Printf("%s was not a known index and will be created with dynamic mapping (no settings)", givenIndex)
			}

			err = createIndex(c, b, givenIndex)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&settings, "settings", "s", "", "Create an index from a json file containing explicit elasticsearch index settings")
}

func createIndex(c *elasticsearch.Client, settings []byte, indexName string) error {
	var (
		res *esapi.Response
		err error

		// Pretty printing of responses
		s            string
		cleanRes     []byte
		m            map[string]any
		marshalErr   error
		unmarshalErr error
	)

	res, err = c.Indices.Delete([]string{indexName})
	if err != nil {
		return fmt.Errorf("cannot delete index: %s", res.Status())
	}

	s = res.String()
	marshalErr = json.Unmarshal([]byte(s), &m)
	cleanRes, unmarshalErr = json.MarshalIndent(m, "", "\t")
	if marshalErr == nil && unmarshalErr == nil {
		s = string(cleanRes)
	}

	if res.IsError() {
		log.Println("warning when deleting index", res.Status())
	} else {
		log.Printf("deleted index %s", indexName)
	}

	res, err = c.Indices.Create(indexName, c.Indices.Create.WithBody(bytes.NewReader(settings)))
	if err != nil {
		return fmt.Errorf("cannot create index: %s", err)
	}

	if res.IsError() {
		marshalErr = json.Unmarshal([]byte(s), &m)
		cleanRes, unmarshalErr = json.MarshalIndent(m, "", "\t")
		if marshalErr == nil && unmarshalErr == nil {
			s = string(cleanRes)
		}
		return fmt.Errorf("cannot create index, got an error response code: %s", s)
	}
	log.Printf("successfully created index %s\n", indexName)
	return nil
}

func isKnownIndex(knownIndices map[string]string, idx string) bool {
	for knownIndex := range knownIndices {
		if idx == knownIndex {
			return true
		}
	}
	return false
}
