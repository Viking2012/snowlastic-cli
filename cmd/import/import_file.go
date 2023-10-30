package _import

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	types "snowlastic-cli/pkg/orm"
)

func GetRecords(fpath string) ([]types.SnowlasticDocument, error) {
	var (
		records []map[string]any
		ret     []types.SnowlasticDocument
		dat     []byte
		err     error
	)

	err = validJSON(fpath)
	if err != nil {
		return nil, err
	}
	log.Println("file passed validation, reading file data...")

	dat, err = os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	log.Printf("%d bytes read, converting to elasticsearch documents...\n", len(dat))
	err = json.Unmarshal(dat, &records)
	if err != nil {
		log.Println("found an error in unmarshalling...")
		return nil, err
	}

	for i := 0; i < len(records); i++ {
		var doc = types.NewDocumentFromMap(records[i])
		ret = append(ret, doc)
	}

	return ret, err
}

func validJSON(fpath string) error {
	var err error
	if ext := path.Ext(fpath); ext != ".json" {
		return fmt.Errorf("file was not a json (%v): %s", ext, fpath)
	}
	_, err = os.Stat(fpath)
	return err
}
