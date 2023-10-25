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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "snowlastic-cli",
	Version: "2.0.0",
	Short:   "Manage, update, and administer an elasticsearch server",
	Long: `Interact with an elasticsearch server, including indexing documents
from a snowflake database. For example:

Create an elasticsearch index from either default settings or a json file.
Index documents from either a snowflake database or a json file.

TODO: Search the elasticsearch index`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setLogLevel()
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		printConfig()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// subcommand additions

	// Flags and configuration settings
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "config file (usually ./snowlastic-cli.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "set verbose output")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().String("settingsDirectory", "", `The settings directory containing elasticsearch settings and document SQL queries`)
	_ = viper.BindPFlag("settingsDirectory", rootCmd.PersistentFlags().Lookup("settingsDirectory"))

	// Snowflake flags
	rootCmd.PersistentFlags().String("snowflakeUser", "", `Snowflake Username`)
	_ = viper.BindPFlag("snowflakeUser", rootCmd.PersistentFlags().Lookup("snowflakeUser"))

	rootCmd.PersistentFlags().String("snowflakePassword", "", `Snowflake Password`)
	_ = viper.BindPFlag("snowflakePassword", rootCmd.PersistentFlags().Lookup("snowflakePassword"))

	rootCmd.PersistentFlags().String("snowflakeAccount", "", `Snowflake Account name`)
	_ = viper.BindPFlag("snowflakeAccount", rootCmd.PersistentFlags().Lookup("snowflakeAccount"))

	rootCmd.PersistentFlags().String("snowflakeWarehouse", "", `Snowflake Warehouse name`)
	_ = viper.BindPFlag("snowflakeWarehouse", rootCmd.PersistentFlags().Lookup("snowflakeWarehouse"))

	rootCmd.PersistentFlags().String("snowflakeRole", "", `Snowflake User role`)
	_ = viper.BindPFlag("snowflakeRole", rootCmd.PersistentFlags().Lookup("snowflakeRole"))

	rootCmd.PersistentFlags().String("snowflakeDatabase", "", `Snowflake Database`)
	_ = viper.BindPFlag("snowflakeDatabase", rootCmd.PersistentFlags().Lookup("snowflakeDatabase"))

	rootCmd.PersistentFlags().StringSlice("snowflakeSchemas", []string{}, "A comma seperated list of relevant schemas")
	_ = viper.BindPFlag("snowflakeSchemas", rootCmd.PersistentFlags().Lookup("snowflakeSchemas"))

	// Elastic flags
	rootCmd.PersistentFlags().String("elasticUrl", "localhost", "URL of the elasticsearch node")
	_ = viper.BindPFlag("elasticUrl", rootCmd.PersistentFlags().Lookup("elasticUrl"))

	rootCmd.PersistentFlags().Int("elasticPort", 9200, "Elasticsearch node port number")
	_ = viper.BindPFlag("elasticPort", rootCmd.PersistentFlags().Lookup("elasticPort"))

	rootCmd.PersistentFlags().String("elasticUser", "", "Elasticsearch username")
	_ = viper.BindPFlag("elasticUser", rootCmd.PersistentFlags().Lookup("elasticUser"))

	rootCmd.PersistentFlags().String("elasticPassword", "", "Elasticsearch user password")
	_ = viper.BindPFlag("elasticPassword", rootCmd.PersistentFlags().Lookup("elasticPassword"))

	rootCmd.PersistentFlags().String("elasticApiKey", "", "Elasticsearch API Key")
	_ = viper.BindPFlag("elasticApiKey", rootCmd.PersistentFlags().Lookup("elasticApiKey"))

	rootCmd.PersistentFlags().String("elasticServiceToken", "", "Elasticsearch Bearer Token")
	_ = viper.BindPFlag("elasticServiceToken", rootCmd.PersistentFlags().Lookup("elasticServiceToken"))

	rootCmd.PersistentFlags().String("elasticCaCertPath", "", "Elasticsearch CA Certificate Path")
	_ = viper.BindPFlag("elasticCaCertPath", rootCmd.PersistentFlags().Lookup("elasticCaCertPath"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".snowlastic-cli" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("snowlastic-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func printConfig() {
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
	fmt.Printf("%-21s: %s\n", "settingsDirectory", viper.GetString("settingsDirectory"))
	fmt.Println()
	fmt.Printf("%-21s: %t\n", "verbose logging", viper.GetBool("verbose"))
	if viper.GetBool("verbose") {
		log.Println("got verbose output")
	}
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

func setLogLevel() {
	log.SetOutput(os.Stderr)
	if viper.GetBool("verbose") {
		log.SetOutput(os.Stdout)
		log.Println("verbose setting received...")
	}
}
