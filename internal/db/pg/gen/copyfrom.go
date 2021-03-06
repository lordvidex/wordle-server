// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: copyfrom.go

package pg

import (
	"context"
)

// iteratorForInsertWords implements pgx.CopyFromSource.
type iteratorForInsertWords struct {
	rows                 []InsertWordsParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertWords) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertWords) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].PlayerGamesID,
		r.rows[0].Word,
		r.rows[0].PlayedAt,
	}, nil
}

func (r iteratorForInsertWords) Err() error {
	return nil
}

func (q *Queries) InsertWords(ctx context.Context, arg []InsertWordsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"player_game_words"}, []string{"player_games_id", "word", "played_at"}, &iteratorForInsertWords{rows: arg})
}
