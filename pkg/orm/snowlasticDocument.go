package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
)

type SnowlasticDocument interface {
	IsDocument()
	GetID() string
	New() SnowlasticDocument
	Get(key string) any
	ScanFrom(*sql.Rows) error
	MarshalJSON() ([]byte, error)
}

type m map[string]any
type Document struct {
	m
}

func (d *Document) IsDocument() {}
func (d *Document) GetID() string {
	var idField = viper.GetString("identifier")
	return fmt.Sprintf("%v", d.m[idField])
}
func (d *Document) New() SnowlasticDocument {
	var m = map[string]any{}
	return &Document{m: m}
}
func (d *Document) Get(key string) any {
	return d.m[key]
}
func (d *Document) ScanFrom(rows *sql.Rows) error {
	var (
		err error

		cols           []string
		columns        []any
		columnPointers []any
	)
	cols, err = rows.Columns()
	if err != nil {
		return err
	}
	columns = make([]any, len(cols))
	columnPointers = make([]any, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}
	// Scan the result into the column pointers...
	if err := rows.Scan(columnPointers...); err != nil {
		return err
	}

	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		d.m[colName] = *val
	}
	return nil
}
func (d *Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.m)
}

func NewDocument() SnowlasticDocument {
	var m = map[string]any{}
	return &Document{m: m}
}
func NewDocumentFromMap(m map[string]any) SnowlasticDocument {
	return &Document{m: m}
}
