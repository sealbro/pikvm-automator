// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"database/sql"
)

type Command struct {
	Created     sql.NullTime
	Updated     sql.NullTime
	ID          string
	Description string
	Expression  string
}
