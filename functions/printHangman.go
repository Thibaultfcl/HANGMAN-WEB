package functions

import (
	"fmt"
	"os"
)

// print the hangman depending of the number of attempts left
func PrintHangman(attemptLeft int) string {
	w := ""
	if attemptLeft != 10 {
		// use the function read to get the right hangman
		if Read("hangman.txt")[9] == byte(13) {
			w = (string(Read("hangman.txt")[(9-attemptLeft)*79 : (9-attemptLeft)*79+77]))
		} else {
			w = (string(Read("hangman.txt")[(9-attemptLeft)*71 : (9-attemptLeft)*71+70]))
		}
	}
	return w
}

// open the file and return it as a byte array
func Read(s string) []byte {
	file, err := os.Open(s)
	if err != nil {
		fmt.Println(err)
	}
	arr := make([]byte, 1200)
	file.Read(arr)
	return arr
}
