-- name: GetUser :one
SELECT id, name, email, created_at FROM users WHERE id = ?;
