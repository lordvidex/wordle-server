--
--  -- QUERIES -- 
--
-- name: GetGameAndSettings :one
SELECT
    *
FROM
    game
    INNER JOIN word ON word.id = game.word_id
    INNER JOIN game_settings ON game_settings.game_id = game.id
WHERE
    game.id = $1
LIMIT
    1;

-- name: ListSessionsInGame :many
SELECT
    *
FROM
    game_session
    INNER JOIN game on game_session.game_id = game.id
    INNER JOIN game_player on game_player.id = game_session.player_id
WHERE
    game_session.game_id = $1;

--
-- CREATIONS 
--
-- name: CreateGame :exec
INSERT INTO
    game (id, word_id)
VALUES
    ($1, $2) RETURNING *;

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

-- name: AddPlayerSessionToGame :exec
INSERT INTO game_session(game_id, player_id) VALUES($1, $2);

-- name: AddPlayerGuess :exec
INSERT INTO game_session_guess(
    game_session_id,
    word_id
) VALUES($1, $2);