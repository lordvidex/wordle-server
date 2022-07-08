-- name: GetGame :one
SELECT * FROM game WHERE id = $1;