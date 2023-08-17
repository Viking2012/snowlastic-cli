package es

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	icm_orm "github.com/alexander-orban/icm_goapi/orm"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/schollz/progressbar/v3"
	"log"
)

const BulkInsertSize = 1000

func BatchEntities(docs <-chan icm_orm.ICMEntity, batchSize int) chan []icm_orm.ICMEntity {
	var batches = make(chan []icm_orm.ICMEntity, 1)

	go func() {
		defer close(batches)
		for keepBatching := true; keepBatching; {
			var batch []icm_orm.ICMEntity
			for {
				select {
				case c, ok := <-docs:
					if !ok {
						keepBatching = false
						goto done
					}
					batch = append(batch, c)
					if len(batch) == batchSize {
						goto done
					}
				}
			}
		done:
			if len(batch) > 0 {
				batches <- batch
			}
		}
	}()

	return batches
}

func BulkImport(es *elasticsearch.Client, batches <-chan []icm_orm.ICMEntity, indexName string, numBatches int64) (numIndexed, numErrors int64, err error) {
	var numProcessed int64 = 1
	bar := progressbar.Default(numBatches)

	for batch := range batches {
		//log.Printf("processing batch #%-5d (%5.1f%%)\n", numProcessed, (float64(numProcessed)/float64(numBatches))*100)
		var buf bytes.Buffer // to collect the bytes of the batch payload
		for _, c := range batch {
			// Prepare the metadata payload
			//
			var idField = c.GetID()
			meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, idField, "\n"))
			data, err := json.Marshal(c)
			if err != nil {
				return numIndexed, numErrors, errors.New(fmt.Sprintf("Cannot encode entity %s: %s", idField, err))
			}
			data = append(data, "\n"...)

			buf.Grow(len(meta) + len(data))
			buf.Write(meta)
			buf.Write(data)

		}
		indexCount, errorCount, err := bulkIndex(es, buf, indexName, BulkInsertSize)
		if err != nil {
			return numIndexed, numErrors, err
		}
		numIndexed += int64(indexCount)
		numErrors += int64(errorCount)
		numProcessed++
		_ = bar.Add(1)
	}

	return numIndexed, numErrors, nil
}

func bulkIndex(es *elasticsearch.Client, buf bytes.Buffer, indexName string, batchSize int) (numIndexed, numErrors int, err error) {
	var (
		res *esapi.Response

		raw map[string]interface{}
		blk *BulkResponse
	)

	res, err = es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex(indexName))
	if err != nil {
		err = errors.New(fmt.Sprintf("Failure indexing batch: %s", err.Error()))
		return numIndexed, numErrors, err
	}
	// If the whole request failed, print error and mark all documents as failed
	//
	if res.IsError() {
		numErrors += batchSize
		if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
			err = errors.New(fmt.Sprintf("Failure to to parse response body: %s", err.Error()))
			return numIndexed, numErrors, err
		} else {
			log.Printf("  Error: [%d] %s: %s",
				res.StatusCode,
				raw["error"].(map[string]interface{})["type"],
				raw["error"].(map[string]interface{})["reason"],
			)
		}
		// A successful response might still contain errors for particular documents...
		//
	} else {
		if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
			log.Fatalf("Failure to to parse response body: %s", err)
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
	_ = res.Body.Close()
	return numIndexed, numErrors, err
}
