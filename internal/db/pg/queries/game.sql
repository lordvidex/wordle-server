--
-- CREATE GAME
--

-- name: CreateGame :exec
INSERT INTO game (id) VALUES ($1) RETURNING *;

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
SELECT * FROM game
    INNER JOIN game_settings ON game_settings.game_id = game.id
    LEFT JOIN word ON word.id = game.word_id
WHERE game.id = $1 LIMIT 1;

-- name: FindByInviteId :many
SELECT * FROM game
         INNER JOIN game_settings gs on game.id = gs.game_id
WHERE invite_id LIKE '%' || $1 || '%';

-- name: GetPlayersInGame :many
SELECT * FROM game_player WHERE game_id = $1;

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