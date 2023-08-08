package demo

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"net/http"
	"os"
	"snowlastic-cli/pkg/es"
	"strconv"
	"strings"
	"time"
)

type Demo struct {
	ID          int    `json:"id"`
	SearchTerm  string `json:"search-term"`
	Value       string `json:"value"`
	ShouldMatch bool   `json:"should-match"`
}

func (d *Demo) IsICMEntity() bool { return true }
func (d *Demo) GetID() string     { return strconv.Itoa(d.ID) }

func IndexDemos(demosPath, demoSettings, credsPath, caCertPath string) error {
	var (
		demos []Demo
		err   error
		b     []byte
		c     *elasticsearch.Client

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

	// Get demos array
	log.Println("reading demos json")
	b, err = os.ReadFile(demosPath)
	if err != nil {

		return errors.New(fmt.Sprintf("error in reading demos json: %s", err))
	}
	log.Println("unmarshalling demos into []Demo")
	err = json.Unmarshal(b, &demos)
	if err != nil {
		return errors.New(fmt.Sprintf("error in unmarshalling demos json: %s", err))
	}

	// create the first index
	log.Println("reading demo settings")
	b, err = os.ReadFile(demoSettings)
	if err != nil {
		return errors.New(fmt.Sprintf("error in reading demos settings: %s", err))
	}

	// get a new ES client
	log.Println("generating new ES client")
	c, err = generateEsClient(credsPath, caCertPath)
	if err != nil {
		return errors.New(fmt.Sprintf("error in generating es client: %s", err))
	}

	log.Printf(
		"\x1b[1mBulk\x1b[0m: documents [%s] batch size [%s]",
		humanize.Comma(int64(len(demos))), humanize.Comma(int64(1000)))
	log.Println(strings.Repeat("▁", 65))

	log.Println("deleting the index, if it exists")
	if res, err = c.Indices.Delete([]string{indexName}); err != nil {
		return errors.New(fmt.Sprintf("Cannot delete index: %s", err))
	}

	log.Println("creating demo index")
	res, err = c.Indices.Create(indexName, c.Indices.Create.WithBody(bytes.NewReader(b)))
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot create index: %s", err))
	}
	if res.IsError() {
		return errors.New(fmt.Sprintf("Cannot create index: %s", res))
	}

	log.Println("creating batches (calculating batch size)")
	if len(demos)%1000 == 0 {
		numBatches = (len(demos) / 1000)
	} else {
		numBatches = (len(demos) / 1000) + 1
	}

	start := time.Now().UTC()

	for i, a := range demos {
		numItems++

		currBatch = i / 1000
		if i == len(demos)-1 {
			currBatch++
		}

		// Prepare the metadata payload
		//
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, a.ID, "\n"))
		// fmt.Printf("%s", meta) // <-- Uncomment to see the payload

		// Prepare the data payload: encode article to JSON
		//
		data, err := json.Marshal(a)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot encode article %d: %s", a.ID, err))
		}

		// Append newline to the data payload
		//
		data = append(data, "\n"...) // <-- Comment out to trigger failure for batch
		// fmt.Printf("%s", data) // <-- Uncomment to see the payload

		// // Uncomment next block to trigger indexing errors -->
		// if a.ID == 11 || a.ID == 101 {
		// 	data = []byte(`{"published" : "INCORRECT"}` + "\n")
		// }
		// // <--------------------------------------------------

		// Append payloads to the buffer (ignoring write errors)
		//
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		// When a threshold is reached, execute the Bulk() request with body from buffer
		//
		if i > 0 && i%1000 == 0 || i == len(demos)-1 {
			fmt.Printf("[%d/%d] ", currBatch, numBatches)

			res, err = c.Bulk(bytes.NewReader(buf.Bytes()), c.Bulk.WithIndex(indexName))
			if err != nil {
				return errors.New(fmt.Sprintf("Failure indexing batch %d: %s", currBatch, err))
			}
			// If the whole request failed, print error and mark all documents as failed
			//
			if res.IsError() {
				numErrors += numItems
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					return errors.New(fmt.Sprintf("Failure to to parse response body: %s", err))
				} else {
					return errors.New(fmt.Sprintf("  Error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					))
				}
				// A successful response might still contain errors for particular documents...
				//
			} else {
				if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
					return errors.New(fmt.Sprintf("Failure to to parse response body: %s", err))
				} else {
					for _, d := range blk.Items {
						// ... so for any HTTP status above 201 ...
						//
						if d.Index.Status > 201 {
							// ... increment the error counter ...
							//
							numErrors++

							// ... and print the response status and error information ...
							log.Printf("  Error: [%d]: %s: %s: %s: %s",
								d.Index.Status,
								d.Index.Error.Type,
								d.Index.Error.Reason,
								d.Index.Error.Cause.Type,
								d.Index.Error.Cause.Reason,
							)
						} else {
							// ... otherwise increase the success counter.
							//
							numIndexed++
						}
					}
				}
			}

			// Close the response body, to prevent reaching the limit for goroutines or file handles
			//
			err = res.Body.Close()
			if err != nil {
				return errors.New(fmt.Sprintf("error in closing the response body: %s", err))
			}

			// Reset the buffer and items counter
			//
			buf.Reset()
			numItems = 0
		}
	}

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	//
	fmt.Print("\n")
	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if numErrors > 0 {
		return errors.New(fmt.Sprintf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			humanize.Comma(int64(numErrors)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		))
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	}
	return nil
}

func generateEsClient(credsPath, caCertPath string) (*elasticsearch.Client, error) {
	// according to https://github.com/elastic/go-elasticsearch/issues/86#issuecomment-527962518
	// required when there are enterprise certificates which need to be used in this context
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	// read the cert used/created by elasticsearch
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return &elasticsearch.Client{}, err
	}
	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		log.Println("No certs appended, using system certs only")
	}

	var creds = make(map[string]string)
	b, err := os.ReadFile(credsPath)
	if err != nil {
		return &elasticsearch.Client{}, err
	}
	err = json.Unmarshal(b, &creds)
	if err != nil {
		return &elasticsearch.Client{}, err
	}

	cfg := elasticsearch.Config{ // works locally, but not when there are enterprise certs
		Addresses: []string{"Https://localhost:9200"},
		Username:  creds["user"],
		Password:  creds["pass"],
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				RootCAs:            rootCAs}},
	}
	c, err := elasticsearch.NewClient(cfg)
	return c, err
}
