-- name: GetCommand :one
SELECT * FROM commands
WHERE id = ? LIMIT 1;

-- name: ListCommands :many
SELECT * FROM commands
ORDER BY id;

-- name: CreateCommand :one
INSERT INTO commands (
    id, description, expression
) VALUES (
             ?, ?, ?
         )
    RETURNING *;

-- name: UpdateCommand :exec
UPDATE commands
set updated = CURRENT_TIMESTAMP,
    description = ?,
    expression = ?
WHERE id = ?;

-- name: DeleteCommand :exec
DELETE FROM commands
WHERE id = ?;