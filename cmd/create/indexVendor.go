package create

import (
	"bytes"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"snowlastic-cli/pkg/es"
)

var (
	err error
	b   []byte
	c   *elasticsearch.Client

	buf bytes.Buffer
	res *esapi.Response
	raw map[string]interface{}
	blk *es.BulkResponse

	indexName = "demo"

	numItems   int
	numErrors  int
	numIndexed int
	numBatches int
	currBatch  int
)

func indexVendor() error { return nil }
