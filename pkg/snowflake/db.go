package snowflake

import (
	"database/sql"
	"github.com/snowflakedb/gosnowflake"
)

type Config struct {
	Account   string
	Warehouse string
	Database  string
	Schema    string
	User      string
	Password  string
	Role      string
}

func NewDB(cfg Config) (*sql.DB, error) {
	var (
		db  *sql.DB
		dsn string
		err error
	)
	var c = gosnowflake.Config{
		Account:   cfg.Account,
		User:      cfg.User,
		Password:  cfg.Password,
		Database:  cfg.Database,
		Schema:    cfg.Schema,
		Warehouse: cfg.Warehouse,
		Role:      cfg.Role,
	}
	dsn, err = gosnowflake.DSN(&c)
	if err != nil {
		return db, err
	}
	db, err = sql.Open("snowflake", dsn)
	return db, err
}
