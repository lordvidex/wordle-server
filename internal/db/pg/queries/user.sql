-- name: GetPlayerByName :many
SELECT * FROM wordlewf_user WHERE name ILIKE '%$1%';

-- name: GetPlayerGames :many
SELECT * from game
     INNER JOIN game_player ON game_player.game_id = game.id
     INNER JOIN wordlewf_user ON wordlewf_user.id = game_player.user_id
     WHERE wordlewf_user.name ILIKE '%' || $1 || '%';