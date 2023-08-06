package snowflake

import (
	"strings"
)

func GenerateSnowflakeDSN(user, password, account, database, schema string) string {
	ret := strings.Builder{}
	ret.WriteString(user)
	ret.WriteString(":")
	ret.WriteString(password)
	ret.WriteString("@")
	ret.WriteString(account)
	if database != "" {
		ret.WriteString("/")
		ret.WriteString(database)
	}
	if schema != "" {
		ret.WriteString("/")
		ret.WriteString(schema)
	}

	return ret.String()
}
