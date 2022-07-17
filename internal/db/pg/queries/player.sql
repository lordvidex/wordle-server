-- name: GetPlayerByEmail :one
SELECT * FROM player WHERE email = $1;

-- name: GetPlayerByID :one
SELECT * FROM player WHERE id = $1;

-- name: GetPlayerByName :many
SELECT * FROM player WHERE name ILIKE '%' || $1 || '%';

-- name: GetPlayerGames :many
SELECT game.*,
     gs.word_length,
     gs.trials,
     gs.max_player_count,
     gs.has_analytics,
     gs.should_record_time,
     gs.can_view_opponents_sessions
 from game
     INNER JOIN game_settings gs on gs.game_id = game.id
     INNER JOIN player_games pg ON game_player.game_id = game.id
     WHERE pg.player_id = $1;

-- name: GetPlayersResultInGame :many
SELECT player.* from player_games pg
     INNER JOIN player ON player.id = pg.player_id
     WHERE pg.game_id = $1;

-- name: GetPlayersResultCountInGame :one
SELECT COUNT(*) from player_games pg WHERE pg.game_id = $1;

-- name: CreatePlayer :one
INSERT INTO player (name, email, password) VALUES ($1, $2, $3) RETURNING *;

-- name: DeletePlayer :exec
UPDATE player SET is_deleted = true WHERE id = $1;