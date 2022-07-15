-- name: GetWord :one
SELECT * from word WHERE id=$1;

-- name: CreateWord :one
INSERT INTO word(
    id,
    time_played, 
    letters
) VALUES ($1, $2, $3) RETURNING *;

-- name: WordsPlayedBy :many
SELECT w.* from game_player_word gpw
         INNER JOIN game_player gp on gpw.player_id = gp.id
         INNER JOIN word w on gpw.word_id = w.id
WHERE gp.id = $1
ORDER BY w.time_played;