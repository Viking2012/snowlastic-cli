package es

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
)

func NewElasticClient(addresses []string, user, pass, apiKey, serviceToken string, caCert []byte) (*elasticsearch.Client, error) {
	// according to https://github.com/elastic/go-elasticsearch/issues/86#issuecomment-527962518
	// required when there are enterprise certificates which need to be used in this context
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		log.Println("No certs appended, using system certs only")
	}

	var cfg elasticsearch.Config = elasticsearch.Config{
		Addresses:    addresses,
		Username:     user,
		Password:     pass,
		APIKey:       apiKey,
		ServiceToken: serviceToken,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				RootCAs:            rootCAs},
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	return es, err
}
