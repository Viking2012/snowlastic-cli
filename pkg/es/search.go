package es

import (
	"encoding/json"
	"strings"
)

//TODO(ajo): consider adding a SearchFieldNested for nested field types.
// Although not used when generating a request, nested types might be included in the response from a search.
// Now that the SearchRequest struct is complete, we now need to research the response to determine
// if a SearchFieldNested struct is requited:
// https://www.elastic.co/guide/en/elasticsearch/reference/8.6/search-search.html#search-api-response-body

const (
	// Field formats
	YYYYMMDD Format = "yyyy-mm-dd"
)

type SearchRequest struct {
	DocValueFields []SearchField      `json:"docvalue_fields"`
	Fields         []SearchField      `json:"fields,omitempty"`
	StoredFields   []string           `json:"stored_fields,omitempty"`
	Explain        bool               `json:"explain,omitempty"`
	From           int                `json:"from,omitempty"`
	IndicesBoost   map[string]float32 `json:"indices_boost,omitempty"`
	MinScore       float32            `json:"min_score,omitempty"`
	Query          Query              `json:"query"`
	Size           int                `json:"size,omitempty"`
	Source         SourceField        `json:"_source,omitempty"`
	Version        bool               `json:"version,omitempty"`
}

func (sr *SearchRequest) MarshalJSON() ([]byte, error) {
	type ret struct {
		DocValueFields []SearchField `json:"docvalue_fields"`
		Fields         []SearchField `json:"fields"`
		StoredFields   string        `json:"stored_fields"`
	}
	return json.Marshal(ret{
		DocValueFields: sr.DocValueFields,
		Fields:         sr.Fields,
		StoredFields:   strings.Join(sr.StoredFields, ","),
	})
}

// SearchField is used in either Docvalue fields or in the source fields, either of which is used
// when elasticsearch generates the response
type SearchField interface {
	IsSearchField() bool
}

type SearchFieldString struct {
	string
}

func (sf *SearchFieldString) IsSearchField() bool { return true }
func (sf *SearchFieldString) MarshalJSON() ([]byte, error) {
	return json.Marshal(sf.string)
}

type SearchFieldObject struct {
	FieldName   string
	FieldFormat Format
}

func (fd *SearchFieldObject) IsSearchField() bool { return true }
func (fd *SearchFieldObject) MarshalJSON() ([]byte, error) {
	var m = make(map[string]interface{})
	m["field"] = fd.FieldName
	m["format"] = &fd.FieldFormat
	return json.Marshal(m)
}

type Format string

func (fd *Format) IsSearchField() bool { return true }
func (fd *Format) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*fd))
}

// Source Field Definitions
type SourceField interface {
	IsSourceField() bool
}

type SourceWildcard string

func (ss *SourceWildcard) IsSourceField() bool { return true }
func (ss *SourceWildcard) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*ss))
}

type SourceBoolean struct {
	bool
}

func (sb *SourceBoolean) IsSourceField() bool { return true }
func (sb *SourceBoolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(sb.bool)
}

type SourceObject struct {
	Include []SourceWildcard `json:"includes,omitempty"`
	Exclude []SourceWildcard `json:"excludes,omitempty"`
}

func (so *SourceObject) IsSourceField() bool { return true }
