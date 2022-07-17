package app

import (
	"fmt"
	"github.com/lordvidex/wordle-wf/internal/adapters"
	"github.com/lordvidex/wordle-wf/internal/words"
	"log"
)

type decorator string 
const (
	bold decorator = "\u001b[1m"
	reset decorator = "\u001b[0m"
)
type color string

const (
	colorGreen  color = "\033[32m"
	colorReset  color = "\033[0m"
	colorGrey   color = "\033[37m"
	colorYellow color = "\033[33m"
)

// bold makes the color bolder
func (c color) bold() string {
	return string(bold) + string(c)
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
	numberOfTries := runIntro()

	stringGenerator := adapters.NewLocalStringGenerator()
	wordsGenerator := words.NewRandomHandler(stringGenerator)
	session := NewSession(numberOfTries, wordsGenerator)
	// start the game loop
	for !session.HasEnded() {
		fmt.Printf("Number of tries used: %d/%d\n", len(session.playedWords), session.maxTries)
		fmt.Println("Enter a word:")
		var word string
		_, err := fmt.Scanln(&word)
		if err != nil {
			log.Fatal("Invalid word", err)
		}
		session.playedWords = append(session.playedWords, words.New(word))
		if session.IsWon() {
			fmt.Println("You won !")
			fmt.Println(colorForStatus(words.Correct).bold(), session.correctWord, colorReset, reset)
			break
		} else {
			fmt.Println("almost there !")
			lastWord := session.playedWords[len(session.playedWords)-1]
			status := lastWord.CompareTo(session.correctWord)
			for i, score := range status {
				fmt.Print(colorForStatus(score).bold(), string(lastWord[i]), colorReset, reset)
			}
			fmt.Print("\n")
		}
	}
	if session.IsWon() {
		fmt.Println("Nice job breaking our code and guessing the word !")
	} else {
		fmt.Println("You tried your best !, the word is ", bold, colorGreen, session.correctWord, colorReset, reset)
	}
	fmt.Println("Would you like to play again? y(Y)/n(N)")
	var input string
	_, err := fmt.Scanln(&input)
	for err != nil {
		fmt.Println("Invalid input, please enter y(Y)/n(N)")
		_, err = fmt.Scanln(&input)
	}
	if input == "y" || input == "Y" {
		Start()
	}
}

type Session struct {
	maxTries    int
	playedWords []words.Word
	correctWord words.Word
}

func NewSession(maxTries int, generator words.RandomHandler) *Session {
	return &Session{
		maxTries:    maxTries,
		correctWord: generator.GetRandomWord(5),
	}
}

func (s *Session) IsWon() bool {
	if len(s.playedWords) == 0 {
		return false
	}
	return s.correctWord == s.playedWords[len(s.playedWords)-1]
}

func (s *Session) HasEnded() bool {
	return s.IsWon() || len(s.playedWords) == s.maxTries
}
