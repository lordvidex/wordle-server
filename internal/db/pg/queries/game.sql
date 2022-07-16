--
-- CREATE GAME
--

-- name: CreateGame :exec
INSERT INTO game (id, invite_id) VALUES ($1, $2) RETURNING *;

-- name: CreateGameSettings :exec
INSERT INTO
    game_settings(
        game_id,
        word_length,
        trials,
        player_count,
        has_analytics,
        should_record_time,
        can_view_opponents_sessions
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

--
-- JOIN GAME
-- 
-- name: CreateGamePlayer :one
INSERT INTO
    game_player (user_id, game_id, name)
VALUES
    ($1, $2, $3) RETURNING *;

--
-- DELETE GAME
--

-- name: DeleteGame :exec
DELETE FROM game WHERE id = $1;

--
-- LEAVE GAME
--

-- name: LeaveGame :exec
-- SOFT DELETE
UPDATE game_player SET deleted = true WHERE id = $1;

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
        player_count,
        has_analytics,
        should_record_time,
        can_view_opponents_sessions
    ) = ($2, $3, $4, $5, $6, $7)
WHERE
    game_settings.game_id = $1 RETURNING *;

--
-- Start Game
-- 
-- name: StartGame :exec
UPDATE game SET start_time = $1 WHERE game.id = $2;

--
-- Play Game (Guess word)
--

-- name: AddGuess :exec
INSERT INTO
    game_player_word(player_id, word_id) VALUES ($1, $2);

-- 
-- Find Game 
-- 

-- name: FindById :one
SELECT game.*,
       game_settings.word_length,
       game_settings.trials,
       game_settings.player_count,
       game_settings.has_analytics,
       game_settings.should_record_time,
       game_settings.can_view_opponents_sessions,
       word.time_played,
       word.letters
       FROM game
    INNER JOIN game_settings ON game_settings.game_id = game.id
    LEFT JOIN word ON word.id = game.word_id
WHERE game.id = $1 LIMIT 1;

-- name: FindByInviteId :many
SELECT game.*,
       gs.word_length,
       gs.trials,
       gs.player_count,
       gs.has_analytics,
       gs.should_record_time,
       gs.can_view_opponents_sessions
       FROM game
         INNER JOIN game_settings gs on game.id = gs.game_id
WHERE
    game.end_time IS NULL -- not ended
  AND
    game.start_time IS NULL -- not started
  AND
    game.invite_id LIKE '%' || $1 || '%'; -- like invite id

-- name: GetPlayersInGame :many
SELECT game_player.*,
       wu.email,
       wu.name as user_name,
       wu.password
FROM game_player
    LEFT JOIN wordlewf_user wu on game_player.user_id = wu.id
WHERE game_id = $1;

--
-- End Game
--

-- name: EndGame :exec
UPDATE
    game
SET
    end_time = $1
WHERE
    game.id = $2;