package es

import (
	"log"
	types "snowlastic-cli/pkg/orm"
	"sync"
	"testing"
)

func TestBatchEntities(t *testing.T) {
	var c = make(chan types.SnowlasticDocument, 1000)
	var want = 200
	var got int

	go func() {
		var wg = sync.WaitGroup{}
		for i := 0; i < want; i++ {
			wg.Add(1)
			var num = i
			go func() {
				var doc = types.NewDocumentFromMap(map[string]any{"": num})
				c <- doc
				wg.Done()
			}()
		}
		wg.Wait()
		close(c)
	}()
	batches := BatchEntities(c, BulkInsertSize)

	var batchNum int
	for batch := range batches {
		log.Printf("processing batch %d", batchNum)
		for range batch {
			got++
		}
		batchNum++
	}
	if got != want {
		t.Errorf("BatchEntities() error: wanted %d documents, but got %d", want, got)
	}
}
