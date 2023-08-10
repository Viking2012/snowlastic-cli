package create

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"snowlastic-cli/pkg/es"
)

func indexDemo() error {
	caCertPath := viper.GetString("elasticCaCertPath")
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return err
	}
	cfg := es.ElasticClientConfig{
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
	c, err := es.NewElasticClient(&cfg)
	if err != nil {
		return err
	}
	res, err := c.Indices.Delete([]string{"demo"})
	if err != nil {
		return fmt.Errorf("cannot delete index: %s", err)
	}
	if res.IsError() {
		log.Println("error when deleting index", res.String())
	} else {
		log.Println(res.String())
	}

	b := []byte(demoIndex)
	res, err = c.Indices.Create("demo", c.Indices.Create.WithBody(bytes.NewReader(b)))
	if err != nil {
		return fmt.Errorf("cannot create index: %s", err)
	}

	if res.IsError() {
		return fmt.Errorf("cannot create index, got an error response code: %s\n", res.String())
	}
	log.Println("successfully created an index")
	return nil
}

const demoIndex string = `{
  "settings": {
    "analysis": {
      "filter": {
        "english_stop": {
          "type":       "stop",
          "stopwords":  "_english_"
        },
        "english_keywords": {
          "type":       "keyword_marker",
          "keywords":   ["example"]
        },
        "english_stemmer": {
          "type":       "stemmer",
          "language":   "english"
        },
        "english_possessive_stemmer": {
          "type":       "stemmer",
          "language":   "possessive_english"
        }
      },
      "analyzer": {
        "demoAnalyzer": {
          "char_filter": ["html_strip"],
          "tokenizer": "standard",
          "filter": [
            "english_possessive_stemmer",
            "english_stop",
            "english_keywords",
            "english_stemmer",
            "lowercase",
            "stop",
            "asciifolding",
            "stemmer",
            "edge_ngram"
          ]
        }
      },
      "normalizer": {
        "demoNormalizer": {
          "type": "custom",
          "filter": ["lowercase","asciifolding"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "search-term":  {"type": "keyword", "normalizer": "demoNormalizer"},
      "value":        {"type": "text"},
      "should-match": {"type": "boolean"}
    }
  }
}
`
