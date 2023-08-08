package create

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func indexDemo(c *elasticsearch.Client) error {
	log.Println("deleting the index, if it exists")
	if res, err = c.Indices.Delete([]string{indexName}); err != nil {
		return errors.New(fmt.Sprintf("Cannot delete index: %s", err))
	}

	log.Println("reading demo settings")
	b := []byte(demoIndex)

	log.Println("creating demo index")
	res, err = c.Indices.Create(indexName, c.Indices.Create.WithBody(bytes.NewReader(b)))
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot create index: %s", err))
	}
	if res.IsError() {
		return errors.New(fmt.Sprintf("Cannot create index: %s", res))
	}

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
