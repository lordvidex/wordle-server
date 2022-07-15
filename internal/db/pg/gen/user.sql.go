// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: user.sql

package pg

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const getPlayerByName = `-- name: GetPlayerByName :many
SELECT id, name, email, password FROM wordlewf_user WHERE name ILIKE '%$1%'
`

func (q *Queries) GetPlayerByName(ctx context.Context) ([]*WordlewfUser, error) {
	rows, err := q.db.Query(ctx, getPlayerByName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*WordlewfUser{}
	for rows.Next() {
		var i WordlewfUser
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
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

const getPlayerGames = `-- name: GetPlayerGames :many
SELECT game.id, invite_id, word_id, start_time, end_time, game_player.id, user_id, game_id, game_player.name, deleted, wordlewf_user.id, wordlewf_user.name, email, password from game
     INNER JOIN game_player ON game_player.game_id = game.id
     INNER JOIN wordlewf_user ON wordlewf_user.id = game_player.user_id
     WHERE wordlewf_user.name ILIKE '%' || $1 || '%'
`

type GetPlayerGamesRow struct {
	ID        uuid.UUID
	InviteID  sql.NullString
	WordID    uuid.NullUUID
	StartTime sql.NullTime
	EndTime   sql.NullTime
	ID_2      uuid.UUID
	UserID    uuid.UUID
	GameID    uuid.UUID
	Name      string
	Deleted   sql.NullBool
	ID_3      uuid.UUID
	Name_2    string
	Email     string
	Password  string
}

func (q *Queries) GetPlayerGames(ctx context.Context, dollar_1 sql.NullString) ([]*GetPlayerGamesRow, error) {
	rows, err := q.db.Query(ctx, getPlayerGames, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetPlayerGamesRow{}
	for rows.Next() {
		var i GetPlayerGamesRow
		if err := rows.Scan(
			&i.ID,
			&i.InviteID,
			&i.WordID,
			&i.StartTime,
			&i.EndTime,
			&i.ID_2,
			&i.UserID,
			&i.GameID,
			&i.Name,
			&i.Deleted,
			&i.ID_3,
			&i.Name_2,
			&i.Email,
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
