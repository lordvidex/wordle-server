-- name: GetWord :one
SELECT * from word WHERE id=@wordId;

-- name: CreateWord :one
INSERT INTO word(
    id,
    time_played, 
    letters
) VALUES (@wordId, @timePlayed, @lettersJSON) RETURNING *;