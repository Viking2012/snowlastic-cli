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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

// CreateCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("create called")
	//},
	//Run: printConfig,
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Snowflake flags
	createCmd.PersistentFlags().String("snowflakeUser", "", `Snowflake Username ("SELECT current_user();")`)
	_ = viper.BindPFlag("snowflakeUser", createCmd.PersistentFlags().Lookup("snowflakeUser"))

	createCmd.PersistentFlags().String("snowflakePassword", "", `Snowflake Password`)
	_ = viper.BindPFlag("snowflakePassword", createCmd.PersistentFlags().Lookup("snowflakePassword"))

	createCmd.PersistentFlags().String("snowflakeAccount", "", `Snowflake Account name (in "https://xyz.us-east-1.azure.snowflakecomputing.com/" then "xyz.us-east-1.azure")`)
	_ = viper.BindPFlag("snowflakeAccount", createCmd.PersistentFlags().Lookup("snowflakeAccount"))

	createCmd.PersistentFlags().String("snowflakeWarehouse", "", `Snowflake Warehouse name ("SELECT current_warehouse();")`)
	_ = viper.BindPFlag("snowflakeWarehouse", createCmd.PersistentFlags().Lookup("snowflakeWarehouse"))

	createCmd.PersistentFlags().String("snowflakeRole", "", `Snowflake User role ("SELECT current_role();")`)
	_ = viper.BindPFlag("snowflakeRole", createCmd.PersistentFlags().Lookup("snowflakeRole"))

	createCmd.PersistentFlags().String("snowflakeDatabase", "", `Snowflake Database (SELECT current_database();)`)
	_ = viper.BindPFlag("snowflakeDatabase", createCmd.PersistentFlags().Lookup("snowflakeDatabase"))

	createCmd.PersistentFlags().StringSlice("snowflakeSchemas", []string{}, "A comma seperated list of relevant schemas")
	_ = viper.BindPFlag("snowflakeSchemas", createCmd.PersistentFlags().Lookup("snowflakeSchemas"))

	// Elastic flags
	createCmd.PersistentFlags().String("elasticUrl", "localhost", "URL of the elasticsearch master node")
	_ = viper.BindPFlag("elasticUrl", createCmd.PersistentFlags().Lookup("elasticUrl"))

	createCmd.PersistentFlags().Int("elasticPort", 9200, "Elasticsearch node port number")
	_ = viper.BindPFlag("elasticPort", createCmd.PersistentFlags().Lookup("elasticPort"))

	createCmd.PersistentFlags().String("elasticUser", "", "Elasticsearch username")
	_ = viper.BindPFlag("elasticUser", createCmd.PersistentFlags().Lookup("elasticUser"))

	createCmd.PersistentFlags().String("elasticPassword", "", "Elasticsearch user password")
	_ = viper.BindPFlag("elasticPassword", createCmd.PersistentFlags().Lookup("elasticPassword"))

	createCmd.PersistentFlags().String("elasticApiKey", "", "Elasticsearch API Kry")
	_ = viper.BindPFlag("elasticApiKey", createCmd.PersistentFlags().Lookup("elasticApiKey"))

	createCmd.PersistentFlags().String("elasticServiceToken", "", "Elasticsearch Bearer Token")
	_ = viper.BindPFlag("elasticServiceToken", createCmd.PersistentFlags().Lookup("elasticServiceToken"))

	createCmd.PersistentFlags().String("elasticCaCertPath", "", "Elasticsearch CA Certificate Path")
	_ = viper.BindPFlag("elasticCaCertPath", createCmd.PersistentFlags().Lookup("elasticCaCertPath"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Add() *cobra.Command {
	return createCmd
}

func printConfig(_ *cobra.Command, _ []string) {
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
	fmt.Printf("%-21s: %s\n", "snowflakeUser", viper.GetString("snowflakeUser"))
	fmt.Printf("%-21s: %s\n", "snowflakePassword", maskedSnowflakePassword.String())
	fmt.Printf("%-21s: %s\n", "snowflakeAccount", viper.GetString("snowflakeAccount"))
	fmt.Printf("%-21s: %s\n", "snowflakeWarehouse", viper.GetString("snowflakeWarehouse"))
	fmt.Printf("%-21s: %s\n", "snowflakeRole", viper.GetString("snowflakeRole"))
	fmt.Printf("%-21s: %s\n", "snowflakeDatabase", viper.GetString("snowflakeDatabase"))
	fmt.Printf("%-21s:", "snowflakeSchemas")
	for i, schema := range viper.GetStringSlice("snowflakeSchemas") {
		if i == 0 {
			fmt.Println(" -", schema)
		} else {
			fmt.Printf("%24s %s\n", "-", schema)
		}
	}
	fmt.Println()
	fmt.Printf("%-21s: %s\n", "elasticUrl", viper.GetString("elasticUrl"))
	fmt.Printf("%-21s: %d\n", "elasticPort", viper.GetInt("elasticPort"))
	fmt.Printf("%-21s: %s\n", "elasticUser", viper.GetString("elasticUser"))
	fmt.Printf("%-21s: %s\n", "elasticPassword", maskedElasticPassword.String())
	fmt.Printf("%-21s: %s\n", "elasticApiKey", viper.GetString("elasticApiKey"))
	fmt.Printf("%-21s: %s\n", "elasticBearerToken", viper.GetString("elasticBearerToken"))
	fmt.Printf("%-21s: %s\n", "elasticCaCertPath", viper.GetString("elasticCaCertPath"))
}
