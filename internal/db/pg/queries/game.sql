--
-- CREATE GAME
--

-- name: CreateGame :exec
INSERT INTO game (id, invite_id, word, player_count, start_time) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: CreateGameSettings :exec
INSERT INTO
    game_settings(
        game_id,
        word_length,
        trials,
        max_player_count,
        has_analytics,
        should_record_time,
        can_view_opponents_sessions
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateGameResult :one
INSERT INTO player_games (player_id, game_id, user_game_name, points, position) VALUES ($1, $2, $3, $4, $5) RETURNING id;

--
-- UPDATE SETTINGS
-- 

-- name: UpdateGameSettings :one
UPDATE
    game_settings
SET
    (
        word_length,
        trials,
        max_player_count,
        has_analytics,
        should_record_time,
        can_view_opponents_sessions
    ) = ($2, $3, $4, $5, $6, $7)
WHERE
    game_settings.game_id = $1 RETURNING *;


-- 
-- Find Game 
-- 

-- name: FindById :one
SELECT game.*,
       game_settings.word_length,
       game_settings.trials,
       game_settings.max_player_count,
       game_settings.has_analytics,
       game_settings.should_record_time,
       game_settings.can_view_opponents_sessions
       FROM game
    INNER JOIN game_settings ON game_settings.game_id = game.id
WHERE game.id = $1 LIMIT 1;

-- name: EndGame :exec
UPDATE
    game
SET
    end_time = $1
WHERE
    game.id = $2;

-- name: DeleteGame :exec
DELETE FROM game WHERE id = $1;