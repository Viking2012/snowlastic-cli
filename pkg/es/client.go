package es

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
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

	var c = elasticsearch.Config{
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
