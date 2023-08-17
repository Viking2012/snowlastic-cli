package create

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

func indexFile(c *elasticsearch.Client, fPath, indexName string) error {
	res, err := c.Indices.Delete([]string{indexName})
	if err != nil {
		return fmt.Errorf("cannot delete index: %s", err)
	}
	if res.IsError() {
		log.Println("error when deleting index", res.String())
	} else {
		log.Println(res.String())
	}

	b, err := os.ReadFile(fPath)
	if err != nil {
		return err
	}
	res, err = c.Indices.Create(indexName, c.Indices.Create.WithBody(bytes.NewReader(b)))
	if err != nil {
		return fmt.Errorf("cannot create index: %s", err)
	}

	if res.IsError() {
		return fmt.Errorf("cannot create index, got an error response code: %s\n", res.String())
	}
	log.Println("successfully created an index")
	return nil
}
