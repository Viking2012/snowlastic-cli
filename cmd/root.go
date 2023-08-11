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
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"snowlastic-cli/cmd/create"
	_import "snowlastic-cli/cmd/import"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "snowlastic-cli",
	Version: "1.0.0",
	Short:   "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetOutput(os.Stdout)
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: printConfig,
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
	rootCmd.AddCommand(create.Add())
	rootCmd.AddCommand(_import.Add())

	// Flags and configuration settings
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "config file (usually ./snowlastic-cli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "set verbose output")
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
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

//func printConfig(_ *cobra.Command, _ []string) {
//	//  Simple print the provided configuration file
//	var (
//		maskedSnowflakePassword strings.Builder
//		maskedElasticPassword   strings.Builder
//	)
//	for range viper.GetString("snowflakePassword") {
//		maskedSnowflakePassword.WriteString("*")
//	}
//	for range viper.GetString("elasticPassword") {
//		maskedElasticPassword.WriteString("*")
//	}
//
//	fmt.Println("------------------------------------------------")
//	fmt.Println("-----  current golastic-cli configuration  -----")
//	fmt.Println("------------------------------------------------")
//	fmt.Println()
//	fmt.Printf("%-21s: %s\n", "snowflakeUser", viper.GetString("snowflakeUser"))
//	fmt.Printf("%-21s: %s\n", "snowflakePassword", maskedSnowflakePassword.String())
//	fmt.Printf("%-21s: %s\n", "snowflakeAccount", viper.GetString("snowflakeAccount"))
//	fmt.Printf("%-21s: %s\n", "snowflakeWarehouse", viper.GetString("snowflakeWarehouse"))
//	fmt.Printf("%-21s: %s\n", "snowflakeRole", viper.GetString("snowflakeRole"))
//	fmt.Printf("%-21s: %s\n", "snowflakeDatabase", viper.GetString("snowflakeDatabase"))
//	fmt.Printf("%-21s:", "snowflakeSchemas")
//	for i, schema := range viper.GetStringSlice("snowflakeSchemas") {
//		if i == 0 {
//			fmt.Println(" -", schema)
//		} else {
//			fmt.Printf("%24s %s\n", "-", schema)
//		}
//	}
//	fmt.Println()
//	fmt.Printf("%-21s: %s\n", "elasticUrl", viper.GetString("elasticUrl"))
//	fmt.Printf("%-21s: %d\n", "elasticPort", viper.GetInt("elasticPort"))
//	fmt.Printf("%-21s: %s\n", "elasticUser", viper.GetString("elasticUser"))
//	fmt.Printf("%-21s: %s\n", "elasticPassword", maskedElasticPassword.String())
//	fmt.Printf("%-21s: %s\n", "elasticApiKey", viper.GetString("elasticApiKey"))
//	fmt.Printf("%-21s: %s\n", "elasticBearerToken", viper.GetString("elasticBearerToken"))
//	fmt.Printf("%-21s: %s\n", "elasticCaCertPath", viper.GetString("elasticCaCertPath"))
//}
