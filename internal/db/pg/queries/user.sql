-- name: GetUserByEmail :one
SELECT * FROM wordlewf_user WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM wordlewf_user WHERE id = $1;

-- name: GetUserByName :many
SELECT * FROM wordlewf_user WHERE name ILIKE '%' || $1 || '%';

-- name: GetUserGames :many
SELECT game.*, game_player.* from game
     INNER JOIN game_player ON game_player.game_id = game.id
     WHERE game_player.user_id = $1;

-- name: InsertUser :one
INSERT INTO wordlewf_user (name, email, password) VALUES ($1, $2, $3) RETURNING *;