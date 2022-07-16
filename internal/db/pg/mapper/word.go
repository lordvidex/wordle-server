package mapper

import (
	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
	"github.com/lordvidex/wordle-wf/internal/words"
)

func WordsPlayedBy(from []*pg.Word) (to []*words.Word, err error) {
	for i, each := range from {
		var letters words.Letters
		err = letters.Scan(each.Letters)
		if err != nil {
			return nil, err
		}
		to[i] = &words.Word{
			ID:         each.ID,
			TimePlayed: each.TimePlayed,
			Letters:    letters,
		}
	}
	return to, nil
}
