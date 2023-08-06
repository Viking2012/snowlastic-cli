package es

import "github.com/elastic/go-elasticsearch/v8"

type Client struct {
	*elasticsearch.TypedClient
}

func NewElasticClient(addresses []string, user, pass, apiKey, serviceToken string, caCert []byte) (*Client, error) {
	var cfg elasticsearch.Config = elasticsearch.Config{
		Addresses:    addresses,
		Username:     user,
		Password:     pass,
		APIKey:       apiKey,
		ServiceToken: serviceToken,
		CACert:       caCert,
	}
	es, err := elasticsearch.NewTypedClient(cfg)
	return &Client{es}, err
}
