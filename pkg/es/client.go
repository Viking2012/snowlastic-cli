package es

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

type ElasticClientConfig struct {
	Addresses    []string
	User         string
	Pass         string
	ApiKey       string
	ServiceToken string
	CaCert       []byte
}

func NewElasticClient(cfg *ElasticClientConfig) (*elasticsearch.Client, error) {
	// according to https://github.com/elastic/go-elasticsearch/issues/86#issuecomment-527962518
	// required when there are enterprise certificates which need to be used in this context
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if ok := rootCAs.AppendCertsFromPEM(cfg.CaCert); !ok {
		log.Println("No certs appended, using system certs only")
	}

	var c elasticsearch.Config = elasticsearch.Config{
		Addresses:    cfg.Addresses,
		Username:     cfg.User,
		Password:     cfg.Pass,
		APIKey:       cfg.ApiKey,
		ServiceToken: cfg.ServiceToken,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				RootCAs:            rootCAs},
		},
	}
	es, err := elasticsearch.NewClient(c)
	return es, err
}

func NewDefaultClient() (*elasticsearch.Client, error) {
	var (
		err        error
		caCert     []byte
		caCertPath string
		cfg        ElasticClientConfig
		c          *elasticsearch.Client
	)

	// generate the CA Certificate bytes needed for the elasticsearch Config
	caCertPath = viper.GetString("elasticCaCertPath")
	caCert, err = os.ReadFile(caCertPath)
	if err != nil {
		return c, err
	}
	cfg = ElasticClientConfig{
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
	return NewElasticClient(&cfg)
}
