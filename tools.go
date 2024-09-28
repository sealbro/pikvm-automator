package pikvm_automator

import (
	_ "embed"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed sql/schema.sql
var SchemaSql string
