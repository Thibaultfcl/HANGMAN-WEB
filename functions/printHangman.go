package functions

import (
	"fmt"
	"os"
)

func PrintHangman(attemptsLeft int) {
	file, err := os.Open("hangman.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	if attemptsLeft >= 1 && attemptsLeft <= 9 {
		PrintLine(file, (9-attemptsLeft)*8, 7)
	}
}
