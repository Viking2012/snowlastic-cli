package cmd

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
	"os"
	"snowlastic-cli/pkg/es"
)

// elasticsearch utilities
func generateElasticClient() (*elasticsearch.Client, error) {
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
