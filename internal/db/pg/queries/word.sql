-- name: InsertWords :copyfrom 
INSERT INTO player_game_words(player_games_id, word, played_at)
VALUES ($1, $2, $3);

-- name: PlayerWordsInGame :many
SELECT pgw.* from player_game_words pgw
INNER JOIN player_games pg ON pgw.player_games_id = pg.id
WHERE pg.player_id=$1 AND pg.game_id=$2;

-- name: WordsPlayedBy :many
SELECT * FROM player_game_words pgw 
INNER JOIN player_games pg ON pgw.player_games_id = pg.id 
WHERE pg.player_id = $1;