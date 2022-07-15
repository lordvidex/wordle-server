// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: game.sql

package pg

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

const addGuess = `-- name: AddGuess :exec

INSERT INTO
    game_player_word(player_id, word_id) VALUES ($1, $2)
`

type AddGuessParams struct {
	PlayerID uuid.UUID
	WordID   uuid.UUID
}

//
// Play Game (Guess word)
//
func (q *Queries) AddGuess(ctx context.Context, arg AddGuessParams) error {
	_, err := q.db.Exec(ctx, addGuess, arg.PlayerID, arg.WordID)
	return err
}

const createGame = `-- name: CreateGame :exec

INSERT INTO game (id, invite_id) VALUES ($1, $2) RETURNING id, invite_id, word_id, start_time, end_time
`

type CreateGameParams struct {
	ID       uuid.UUID
	InviteID string
}

//
// CREATE GAME
//
func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) error {
	_, err := q.db.Exec(ctx, createGame, arg.ID, arg.InviteID)
	return err
}

const createGamePlayer = `-- name: CreateGamePlayer :one
INSERT INTO
    game_player (user_id, game_id, name)
VALUES
    ($1, $2, $3) RETURNING id, user_id, game_id, name, deleted
`

type CreateGamePlayerParams struct {
	UserID uuid.UUID
	GameID uuid.UUID
	Name   string
}

//
// JOIN GAME
//
func (q *Queries) CreateGamePlayer(ctx context.Context, arg CreateGamePlayerParams) (*GamePlayer, error) {
	row := q.db.QueryRow(ctx, createGamePlayer, arg.UserID, arg.GameID, arg.Name)
	var i GamePlayer
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Deleted,
	)
	return &i, err
}

const createGameSettings = `-- name: CreateGameSettings :exec
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
    ($1, $2, $3, $4, $5, $6, $7) RETURNING id, game_id, word_length, trials, player_count, has_analytics, should_record_time, can_view_opponents_sessions
`

type CreateGameSettingsParams struct {
	GameID                   uuid.NullUUID
	WordLength               sql.NullInt16
	Trials                   sql.NullInt16
	PlayerCount              sql.NullInt16
	HasAnalytics             sql.NullBool
	ShouldRecordTime         sql.NullBool
	CanViewOpponentsSessions sql.NullBool
}

func (q *Queries) CreateGameSettings(ctx context.Context, arg CreateGameSettingsParams) error {
	_, err := q.db.Exec(ctx, createGameSettings,
		arg.GameID,
		arg.WordLength,
		arg.Trials,
		arg.PlayerCount,
		arg.HasAnalytics,
		arg.ShouldRecordTime,
		arg.CanViewOpponentsSessions,
	)
	return err
}

const deleteGame = `-- name: DeleteGame :exec

DELETE FROM game WHERE id = $1
`

//
// DELETE GAME
//
func (q *Queries) DeleteGame(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteGame, id)
	return err
}

const endGame = `-- name: EndGame :exec

UPDATE
    game
SET
    end_time = $1
WHERE
    game.id = $2
`

type EndGameParams struct {
	EndTime sql.NullTime
	ID      uuid.UUID
}

//
// End Game
//
func (q *Queries) EndGame(ctx context.Context, arg EndGameParams) error {
	_, err := q.db.Exec(ctx, endGame, arg.EndTime, arg.ID)
	return err
}

const findById = `-- name: FindById :one

SELECT game.id, game.invite_id, game.word_id, game.start_time, game.end_time,
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
WHERE game.id = $1 LIMIT 1
`

type FindByIdRow struct {
	ID                       uuid.UUID
	InviteID                 string
	WordID                   uuid.NullUUID
	StartTime                sql.NullTime
	EndTime                  sql.NullTime
	WordLength               sql.NullInt16
	Trials                   sql.NullInt16
	PlayerCount              sql.NullInt16
	HasAnalytics             sql.NullBool
	ShouldRecordTime         sql.NullBool
	CanViewOpponentsSessions sql.NullBool
	TimePlayed               sql.NullTime
	Letters                  pgtype.JSON
}

//
// Find Game
//
func (q *Queries) FindById(ctx context.Context, id uuid.UUID) (*FindByIdRow, error) {
	row := q.db.QueryRow(ctx, findById, id)
	var i FindByIdRow
	err := row.Scan(
		&i.ID,
		&i.InviteID,
		&i.WordID,
		&i.StartTime,
		&i.EndTime,
		&i.WordLength,
		&i.Trials,
		&i.PlayerCount,
		&i.HasAnalytics,
		&i.ShouldRecordTime,
		&i.CanViewOpponentsSessions,
		&i.TimePlayed,
		&i.Letters,
	)
	return &i, err
}

const findByInviteId = `-- name: FindByInviteId :many
SELECT game.id, invite_id, word_id, start_time, end_time, gs.id, game_id, word_length, trials, player_count, has_analytics, should_record_time, can_view_opponents_sessions FROM game
         INNER JOIN game_settings gs on game.id = gs.game_id
WHERE invite_id LIKE '%' || $1 || '%'
`

type FindByInviteIdRow struct {
	ID                       uuid.UUID
	InviteID                 string
	WordID                   uuid.NullUUID
	StartTime                sql.NullTime
	EndTime                  sql.NullTime
	ID_2                     uuid.UUID
	GameID                   uuid.NullUUID
	WordLength               sql.NullInt16
	Trials                   sql.NullInt16
	PlayerCount              sql.NullInt16
	HasAnalytics             sql.NullBool
	ShouldRecordTime         sql.NullBool
	CanViewOpponentsSessions sql.NullBool
}

func (q *Queries) FindByInviteId(ctx context.Context, dollar_1 sql.NullString) ([]*FindByInviteIdRow, error) {
	rows, err := q.db.Query(ctx, findByInviteId, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*FindByInviteIdRow{}
	for rows.Next() {
		var i FindByInviteIdRow
		if err := rows.Scan(
			&i.ID,
			&i.InviteID,
			&i.WordID,
			&i.StartTime,
			&i.EndTime,
			&i.ID_2,
			&i.GameID,
			&i.WordLength,
			&i.Trials,
			&i.PlayerCount,
			&i.HasAnalytics,
			&i.ShouldRecordTime,
			&i.CanViewOpponentsSessions,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPlayersInGame = `-- name: GetPlayersInGame :many
SELECT game_player.id, game_player.user_id, game_player.game_id, game_player.name, game_player.deleted,
       wu.email,
       wu.name as user_name,
       wu.password
FROM game_player
    LEFT JOIN wordlewf_user wu on game_player.user_id = wu.id
WHERE game_id = $1
`

type GetPlayersInGameRow struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	GameID   uuid.UUID
	Name     string
	Deleted  sql.NullBool
	Email    sql.NullString
	UserName sql.NullString
	Password sql.NullString
}

func (q *Queries) GetPlayersInGame(ctx context.Context, gameID uuid.UUID) ([]*GetPlayersInGameRow, error) {
	rows, err := q.db.Query(ctx, getPlayersInGame, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetPlayersInGameRow{}
	for rows.Next() {
		var i GetPlayersInGameRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.GameID,
			&i.Name,
			&i.Deleted,
			&i.Email,
			&i.UserName,
			&i.Password,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const leaveGame = `-- name: LeaveGame :exec

UPDATE game_player SET deleted = true WHERE id = $1
`

//
// LEAVE GAME
//
// SOFT DELETE
func (q *Queries) LeaveGame(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, leaveGame, id)
	return err
}

const startGame = `-- name: StartGame :exec
UPDATE game SET start_time = $1 WHERE game.id = $2
`

type StartGameParams struct {
	StartTime sql.NullTime
	ID        uuid.UUID
}

//
// Start Game
//
func (q *Queries) StartGame(ctx context.Context, arg StartGameParams) error {
	_, err := q.db.Exec(ctx, startGame, arg.StartTime, arg.ID)
	return err
}

const updateGameSettings = `-- name: UpdateGameSettings :one
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
    game_settings.game_id = $1 RETURNING id, game_id, word_length, trials, player_count, has_analytics, should_record_time, can_view_opponents_sessions
`

type UpdateGameSettingsParams struct {
	GameID                   uuid.NullUUID
	WordLength               sql.NullInt16
	Trials                   sql.NullInt16
	PlayerCount              sql.NullInt16
	HasAnalytics             sql.NullBool
	ShouldRecordTime         sql.NullBool
	CanViewOpponentsSessions sql.NullBool
}

//
// UPDATE SETTINGS
//
func (q *Queries) UpdateGameSettings(ctx context.Context, arg UpdateGameSettingsParams) (*GameSetting, error) {
	row := q.db.QueryRow(ctx, updateGameSettings,
		arg.GameID,
		arg.WordLength,
		arg.Trials,
		arg.PlayerCount,
		arg.HasAnalytics,
		arg.ShouldRecordTime,
		arg.CanViewOpponentsSessions,
	)
	var i GameSetting
	err := row.Scan(
		&i.ID,
		&i.GameID,
		&i.WordLength,
		&i.Trials,
		&i.PlayerCount,
		&i.HasAnalytics,
		&i.ShouldRecordTime,
		&i.CanViewOpponentsSessions,
	)
	return &i, err
}
