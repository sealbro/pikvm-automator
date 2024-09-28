package sqlite

import (
	"context"
	"database/sql"
)

func New(fileName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ApplySchema(ctx context.Context, db *sql.DB, ddl string) error {
	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	return nil
}
