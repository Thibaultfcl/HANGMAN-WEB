package functions

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

//get a random world based on the chosen difficulty
func Random(difficulty string) string {
	file := ""
	word := ""

	// get a random index based on the difficulty
	var index int
	if difficulty == "EASY" {
		index = rand.Intn(37) + 1
		file = "words.txt"
	} else if difficulty == "NORMAL" {
		index = rand.Intn(23) + 1
		file = "words2.txt"
	} else if difficulty == "HARD" {
		index = rand.Intn(24) + 1
		file = "words3.txt"
	} else {
		fmt.Println("Erreur")
		os.Exit(1)
	}

	// open the file
	fileOpen, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	i := 0
	// scan it to get the word corresponding to the index
	scanner := bufio.NewScanner(fileOpen)
	for scanner.Scan() {
		i++
		if i == index {
			word = scanner.Text()
		}
	}
	return word
}
