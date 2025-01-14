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
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
	"snowlastic-cli/pkg/es"
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
	Version: "3.1.1",
	Short:   "Manage, update, and administer an elasticsearch server",
	Long: `Interact with an elasticsearch server, including indexing documents
from a snowflake database. For example:

Create an elasticsearch index from either default settings or a json file.
Index documents from either a snowflake database or a json file.

TODO: Search the elasticsearch index`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		setLogLevel()
		log.Println("setting proxies")
		err = setProxy()
		if err != nil {
			return err
		}
		return nil
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

	var indexMap map[string]map[string]string
	//indexMap, err := convertIndices(rootCmd.PersistentFlags().Lookup("elasticIndices"))
	err := viper.UnmarshalKey("elasticIndices", &indexMap)
	if err != nil {
		panic(err)
	}
	viper.Set("elasticIndices", indexMap)
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
	fmt.Println()
	fmt.Printf("%-21s: %s\n", "elasticUrl", viper.GetString("elasticUrl"))
	fmt.Printf("%-21s: %d\n", "elasticPort", viper.GetInt("elasticPort"))
	fmt.Printf("%-21s: %s\n", "elasticUser", viper.GetString("elasticUser"))
	fmt.Printf("%-21s: %s\n", "elasticPassword", maskedElasticPassword.String())
	fmt.Printf("%-21s: %s\n", "elasticApiKey", viper.GetString("elasticApiKey"))
	fmt.Printf("%-21s: %s\n", "elasticBearerToken", viper.GetString("elasticBearerToken"))
	fmt.Printf("%-21s: %s\n", "elasticCaCertPath", viper.GetString("elasticCaCertPath"))
	fmt.Println()

	for k, v := range viper.Get("elasticIndices").(map[string]map[string]string) {
		fmt.Printf("%24s: %24s\n", k, "")
		for k2, v2 := range v {
			fmt.Printf("%24s %-24s:%s\n", "", k2, v2)
		}
	}
}

func setLogLevel() {
	log.SetOutput(os.Stderr)
	if viper.GetBool("verbose") {
		log.SetOutput(os.Stdout)
		log.Println("verbose setting received...")
	}
}

// elasticsearch utilities
func generateDefaultElasticClient() (*elasticsearch.Client, error) {
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
		return c, err
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
	return es.NewElasticClient(&cfg)
}

func setProxy() error {
	var err error
	proxy := getProxy()
	err = os.Setenv(`HTTP_PROXY`, proxy)
	if err != nil {
		return err
	}
	err = os.Setenv(`HTTPS_PROXY`, proxy)
	if err != nil {
		return err
	}
	return nil
}

func getProxy() string {
	proxy := os.Getenv(`ELT_PROXY`)
	if len(proxy) > 7 && proxy[:7] != `http://` {
		proxy = fmt.Sprintf(`http://%v`, proxy)
	}
	return proxy
}
