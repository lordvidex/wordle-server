package remote

import (
	"database/sql"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/lordvidex/wordle-wf/cmd/client/local"
	"github.com/lordvidex/wordle-wf/internal/adapters"
	"github.com/lordvidex/wordle-wf/internal/words"
	"github.com/rivo/tview"
	"os"
)

var (
	app = tview.NewApplication()
)

type WordleCell struct {
	letter       rune
	letterStatus words.LetterStatus
	Empty        bool
}

func (c WordleCell) TviewCell() *tview.TableCell {
	var t string
	if c.Empty {
		t = " "
	} else {
		t = string(c.letter)
	}
	return &tview.TableCell{
		Text:            t,
		Align:           tview.AlignCenter,
		MaxWidth:        0,
		Color:           c.textColor(),
		BackgroundColor: c.backgroundColor(),
		Transparent:     false,
		Attributes:      tcell.AttrBold,
		NotSelectable:   true,
		Clicked:         nil,
	}
}

func (c WordleCell) textColor() tcell.Color {
	if c.letterStatus == words.Unknown {
		return tcell.ColorBlack
	}
	return tcell.ColorWhite
}

func (c WordleCell) backgroundColor() tcell.Color {
	if c.Empty {
		return tcell.ColorWhite
	}
	switch c.letterStatus {
	case words.Correct:
		return tcell.ColorGreen
	case words.Incorrect:
		return tcell.ColorDarkGray
	case words.Exists:
		return tcell.ColorLightGoldenrodYellow
	default:
		return tcell.ColorWhite
	}
}

type WordleBoard struct {
	tview.TableContentReadOnly
	Session     *local.Session
	rowCount    int
	columnCount int
	activeRow   int
	activeWord  []rune
	letterChan  chan rune
	sigChan     chan tcell.Key
	done        chan interface{}
}

func NewBoard(wordLength int, session *local.Session, done chan interface{}) *WordleBoard {
	board := &WordleBoard{
		rowCount:    session.MaxTries,
		columnCount: wordLength,
		letterChan:  make(chan rune, 1),
		sigChan:     make(chan tcell.Key, 1),
		done:        done,
		Session:     session,
	}
	go func() {
		for {
			select {
			case r := <-board.letterChan:
				if len(board.activeWord) < wordLength {
					board.activeWord = append(board.activeWord, r)
				}
			case sig := <-board.sigChan:
				switch sig {
				case tcell.KeyEnter:
					if len(board.activeWord) == wordLength {
						board.EnterWord()
						board.Judge()
					}
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					if len(board.activeWord) > 0 {
						board.activeWord = board.activeWord[:len(board.activeWord)-1]
					}
				}

			case <-done:
				close(board.letterChan)
				close(board.sigChan)
				return
			}

		}
	}()
	return board
}

func (w *WordleBoard) GetCell(row, column int) *tview.TableCell {
	var cell WordleCell
	if row < len(w.Session.PlayedWords) {
		word := w.Session.PlayedWords[row]
		cell = WordleCell{
			letter:       word.Runes()[column],
			letterStatus: word.Stats[column],
		}
	} else {
		if row == w.activeRow && column < len(w.activeWord) {
			cell = WordleCell{
				letter:       w.activeWord[column],
				letterStatus: words.Unknown,
			}
		} else {
			cell = WordleCell{Empty: true}
		}
	}
	return cell.TviewCell()
}

func (w *WordleBoard) EnterWord() {
	status := w.Session.WordStatus(words.New(string(w.activeWord)))
	w.Session.PlayedWords = append(w.Session.PlayedWords, words.Word{
		Word:     string(w.activeWord),
		PlayedAt: sql.NullTime{Valid: false},
		Stats:    status,
	})
	w.activeWord = []rune{}
	w.activeRow++
}

func (w *WordleBoard) Judge() {
	if w.Session.IsWon() {
		fmt.Println("correct word")
	}
	if w.Session.HasEnded() {
		app.Stop()
	}
}

func (w *WordleBoard) GetRowCount() int {
	return w.rowCount
}

func (w *WordleBoard) GetColumnCount() int {
	return w.columnCount
}

func buildTable(col int, session *local.Session, done chan interface{}) *tview.Table {
	board := NewBoard(col, session, done)
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(false, false).
		SetContent(board)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// process text input
		r := event.Rune()
		var text string
		switch {
		case r >= 'A' && r <= 'Z':
			text = string(r)
		case r >= 'a' && r <= 'z':
			text = string(r - 'a' + 'A') // to uppercase
		}
		if text != "" {
			board.letterChan <- rune(text[0])
		} else {
			// process key input
			board.sigChan <- event.Key()
		}

		return event
	})
	return table
}

func Start() {
	done := make(chan interface{}, 1)
	stringGenerator := adapters.NewLocalStringGenerator()
	wordsGenerator := words.NewRandomHandler(stringGenerator)
	session := local.NewSession(5, wordsGenerator)
	table := buildTable(5, session, done)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			done <- "quit"
		}
		return event
	})
	if err := app.SetRoot(table, false).
		Run(); err != nil {
		os.Exit(1)
	}
}
