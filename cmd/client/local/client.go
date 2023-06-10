package local

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/lordvidex/wordle-wf/internal/words"
)

type decorator string

const (
	bold  decorator = "\u001b[1m"
	reset decorator = "\u001b[0m"
)

func (d decorator) Decorate(text string) string {
	return string(d) + text + string(reset)
}

type color string

const (
	colorGreen  color = "\033[32m"
	colorReset  color = "\033[0m"
	colorGrey   color = "\033[37m"
	colorYellow color = "\033[33m"
)

func (c color) Colored(text string) string {
	return string(c) + text + string(colorReset)
}

func colorForStatus(status words.LetterStatus) color {
	switch status {
	case words.Correct:
		return colorGreen
	case words.Exists:
		return colorYellow
	default:
		return colorGrey
	}
}

// runIntro introduces the game and lets us know how many tries the user want before he loses the game
func runIntro() int {
	// introduce the game
	fmt.Printf("Welcome to %sW%sO%sR%sD%sL%sE%s Local !\n", colorGreen, colorYellow, colorYellow, colorGreen, colorGrey, colorGreen, colorReset)
	return getNumberOfTries()
}

func getNumberOfTries() int {
	fmt.Println("How many tries do you want to have per round?")
	var numberOfTries int
	_, err := fmt.Scanln(&numberOfTries)
	if err != nil {
		log.Fatal("Invalid number of tries ", err)
	}
	return numberOfTries
}

func Start() {
	if _, err := tea.NewProgram(NewIntroModel(), tea.WithMouseCellMotion()).Run(); err != nil {
		log.Println("error running the game", err)
	}
	// numberOfTries := runIntro()
	//
	// stringGenerator := adapters.NewLocalStringGenerator()
	// wordsGenerator := words.NewRandomHandler(stringGenerator)
	// session := NewSession(numberOfTries, wordsGenerator)
	// // start the game loop
	// for !session.HasEnded() {
	// 	fmt.Printf("Number of tries used: %d/%d\n", len(session.PlayedWords), session.MaxTries)
	// 	fmt.Println("Enter a word:")
	// 	var word string
	// 	_, err := fmt.Scanln(&word)
	// 	if err != nil {
	// 		log.Fatal("Invalid word", err)
	// 	}
	// 	session.PlayedWords = append(session.PlayedWords, words.New(word))
	// 	if session.IsWon() {
	// 		fmt.Println("You won !")
	// 		fmt.Println(bold.Decorate(colorForStatus(words.Correct).Colored(session.correctWord.String())))
	// 		break
	// 	} else {
	// 		fmt.Println("almost there !")
	// 		lastWord := session.PlayedWords[len(session.PlayedWords)-1]
	// 		status := session.WordStatus(lastWord)
	// 		for i, score := range status {
	// 			fmt.Print(bold.Decorate(colorForStatus(score).Colored(string(lastWord.Word[i]))))
	// 		}
	// 		fmt.Print("\n")
	// 	}
	// }
	// if session.IsWon() {
	// 	fmt.Println("Nice job breaking our code and guessing the word !")
	// } else {
	// 	fmt.Println("You tried your best !, the word is ", bold.Decorate(colorGreen.Colored(session.correctWord.String())))
	// }
	// fmt.Println("Would you like to play again? y(Y)/n(N)")
	// var input string
	// _, err := fmt.Scanln(&input)
	// for err != nil {
	// 	fmt.Println("Invalid input, please enter y(Y)/n(N)")
	// 	_, err = fmt.Scanln(&input)
	// }
	// if input == "y" || input == "Y" {
	// 	Start()
	// }
}

type Session struct {
	MaxTries    int
	PlayedWords []words.Word
	correctWord words.Word
}

func NewSession(maxTries int, generator words.RandomHandler) *Session {
	return &Session{
		MaxTries:    maxTries,
		correctWord: generator.GetRandomWord(5),
	}
}

func (s *Session) WordStatus(word words.Word) []words.LetterStatus {
	return word.CompareTo(s.correctWord)
}

func (s *Session) IsWon() bool {
	if len(s.PlayedWords) == 0 {
		return false
	}
	return s.correctWord.String() == s.PlayedWords[len(s.PlayedWords)-1].String()
}

func (s *Session) HasEnded() bool {
	return s.IsWon() || len(s.PlayedWords) == s.MaxTries
}
