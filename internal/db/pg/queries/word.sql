-- name: GetWord :one
SELECT * from word WHERE id=$1;

-- name: CreateWord :one
INSERT INTO word(
    id,
    time_played, 
    letters
) VALUES ($1, $2, $3) RETURNING *;