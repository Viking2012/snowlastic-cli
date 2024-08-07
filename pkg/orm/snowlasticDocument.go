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
	for k := range result {
		v := result[k]
		w := deepConvertStrings(v)
		result[k] = w
	}
	d.m = keysToLower(result)
	//d.m = result
	return nil
}
func (d *Document) MarshalJSON() ([]byte, error) {
	var m = make(map[string]any)
	for k, v := range d.m {
		if v != nil && k != viper.GetString("identifier") {
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

func deepConvertStrings(v any) any {
	if v == nil {
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return v
	}
	var err error
	var m = make(map[string]any)
	err = json.Unmarshal([]byte(s), &m)
	if err == nil {
		for key := range m {
			m[key] = deepConvertStrings(m[key])
		}
		return m
	}
	var n []map[string]any
	err = json.Unmarshal([]byte(s), &n)
	if err == nil {
		for i := range n {
			var p = n[i]
			for key := range p {
				p[key] = deepConvertStrings(p[key])
			}
			n[i] = p
		}
		return n
	}
	return v
}
