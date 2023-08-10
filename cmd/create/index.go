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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var (
	isVendor   bool
	isCustomer bool
	isDemo     bool
	runAll     bool
	fromFile   string
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
		printConfigFromIndex(cmd, args)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
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
			err = indexDemo()
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
			err = indexDemo()
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

func indexCustomer() error { return nil }

func indexFromFile(fp string) error { return nil }

func printConfigFromIndex(_ *cobra.Command, _ []string) {
	//  Simple print the provided configuration file
	var (
		maskedSnowflakePassword strings.Builder
		maskedElasticPassword   strings.Builder
	)
	for range viper.GetString("snowflakePassword") {
		maskedSnowflakePassword.WriteString("*")
	}
	for range viper.GetString("elasticPassword") {
		maskedElasticPassword.WriteString("*")
	}

	fmt.Println("------------------------------------------------")
	fmt.Println("-----  current golastic-cli configuration  -----")
	fmt.Println("------------------------------------------------")
	fmt.Println()
	fmt.Printf("%-21s: '%s'\n", "snowflakeUser", viper.GetString("snowflakeUser"))
	fmt.Printf("%-21s: '%s'\n", "snowflakePassword", viper.GetString("snowflakePassword"))
	fmt.Printf("%-21s: '%s'\n", "snowflakeAccount", viper.GetString("snowflakeAccount"))
	fmt.Printf("%-21s: '%s'\n", "snowflakeWarehouse", viper.GetString("snowflakeWarehouse"))
	fmt.Printf("%-21s: '%s'\n", "snowflakeRole", viper.GetString("snowflakeRole"))
	fmt.Printf("%-21s: '%s'\n", "snowflakeDatabase", viper.GetString("snowflakeDatabase"))
	fmt.Printf("%-21s:", "snowflakeSchemas")
	for i, schema := range viper.GetStringSlice("snowflakeSchemas") {
		if i == 0 {
			fmt.Println(" -", schema)
		} else {
			fmt.Printf("%24s '%s'\n", "-", schema)
		}
	}
	fmt.Println()
	fmt.Printf("%-21s: '%s'\n", "elasticUrl", viper.GetString("elasticUrl"))
	fmt.Printf("%-21s: '%d'\n", "elasticPort", viper.GetInt("elasticPort"))
	fmt.Printf("%-21s: '%s'\n", "elasticUser", viper.GetString("elasticUser"))
	fmt.Printf("%-21s: '%s'\n", "elasticPassword", viper.GetString("elasticPassword"))
	fmt.Printf("%-21s: '%s'\n", "elasticApiKey", viper.GetString("elasticApiKey"))
	fmt.Printf("%-21s: '%s'\n", "elasticBearerToken", viper.GetString("elasticBearerToken"))
	fmt.Printf("%-21s: '%s'\n", "elasticCaCertPath", viper.GetString("elasticCaCertPath"))
}
