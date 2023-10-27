package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"strings"
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

		result = make(map[string]any)
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
		result[colName] = *val
	}
	d.m = keysToLower(result)
	//d.m = result
	return nil
}
func (d *Document) MarshalJSON() ([]byte, error) {
	var m = make(map[string]any)
	for k, v := range d.m {
		if v != nil {
			m[k] = v
		}
	}
	return json.Marshal(m)
	//return json.Marshal(d.m)
}

func NewDocument() SnowlasticDocument {
	var m = map[string]any{}
	return &Document{m: m}
}
func NewDocumentFromMap(m map[string]any) SnowlasticDocument {
	return &Document{m: keysToLower(m)}
	//return &Document{m: m}
}

func keysToLower(m map[string]any) map[string]any {
	var n = make(map[string]any)
	for key := range m {
		lowercaseKey := strings.ToLower(key)
		var val = m[key]
		n[lowercaseKey] = lower(val)
	}
	return n
}

func lower(v any) any {
	switch v := v.(type) {
	case []any:
		lv := make([]any, len(v))
		for i := range v {
			lv[i] = lower(v[i])
		}
		return lv
	case map[string]any:
		lv := make(map[string]any, len(v))
		for mk, mv := range v {
			lv[strings.ToLower(mk)] = lower(mv)
		}
		return lv
	default:
		return v
	}
}
