// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: commands_query.sql

package repository

import (
	"context"
)

const createCommand = `-- name: CreateCommand :one
INSERT INTO commands (
    id, description, expression
) VALUES (
             ?, ?, ?
         )
    RETURNING created, updated, id, description, expression
`

type CreateCommandParams struct {
	ID          string
	Description string
	Expression  string
}

func (q *Queries) CreateCommand(ctx context.Context, arg CreateCommandParams) (Command, error) {
	row := q.db.QueryRowContext(ctx, createCommand, arg.ID, arg.Description, arg.Expression)
	var i Command
	err := row.Scan(
		&i.Created,
		&i.Updated,
		&i.ID,
		&i.Description,
		&i.Expression,
	)
	return i, err
}

const deleteCommand = `-- name: DeleteCommand :exec
DELETE FROM commands
WHERE id = ?
`

func (q *Queries) DeleteCommand(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteCommand, id)
	return err
}

const getCommand = `-- name: GetCommand :one
SELECT created, updated, id, description, expression FROM commands
WHERE id = ? LIMIT 1
`

func (q *Queries) GetCommand(ctx context.Context, id string) (Command, error) {
	row := q.db.QueryRowContext(ctx, getCommand, id)
	var i Command
	err := row.Scan(
		&i.Created,
		&i.Updated,
		&i.ID,
		&i.Description,
		&i.Expression,
	)
	return i, err
}

const listCommands = `-- name: ListCommands :many
SELECT created, updated, id, description, expression FROM commands
ORDER BY id
`

func (q *Queries) ListCommands(ctx context.Context) ([]Command, error) {
	rows, err := q.db.QueryContext(ctx, listCommands)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Command
	for rows.Next() {
		var i Command
		if err := rows.Scan(
			&i.Created,
			&i.Updated,
			&i.ID,
			&i.Description,
			&i.Expression,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCommand = `-- name: UpdateCommand :exec
UPDATE commands
set updated = CURRENT_TIMESTAMP,
    description = ?,
    expression = ?
WHERE id = ?
`

type UpdateCommandParams struct {
	Description string
	Expression  string
	ID          string
}

func (q *Queries) UpdateCommand(ctx context.Context, arg UpdateCommandParams) error {
	_, err := q.db.ExecContext(ctx, updateCommand, arg.Description, arg.Expression, arg.ID)
	return err
}