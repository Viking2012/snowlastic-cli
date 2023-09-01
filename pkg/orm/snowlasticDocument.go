package types

import "database/sql"

type SnowlasticDocument interface {
	IsDocument()
	GetID() string
	GetQuery(string, string) string
	New() SnowlasticDocument
	ScanFrom(*sql.Rows) error
}
