package mapper

import (
	"database/sql"

	"github.com/google/uuid"
	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
	"github.com/lordvidex/wordle-wf/internal/words"
)

func InsertWords(playerGamesID uuid.NullUUID, w []*words.Word) []pg.InsertWordsParams {
	result := make([]pg.InsertWordsParams, len(w))
	for i, word := range w {
		result[i] = pg.InsertWordsParams{
			Word:          word.String(),
			PlayerGamesID: playerGamesID,
			PlayedAt:      word.PlayedAt.Time,
		}
	}
	return result
}

func PlayerWordsInGame(list []*pg.PlayerGameWord) []*words.Word {
	ans := make([]*words.Word, len(list))
	for i, w := range list {
		ans[i] = &words.Word{
			Word: w.Word,
			PlayedAt: sql.NullTime{Time: w.PlayedAt, Valid: true},
		}
	}
	return ans
}