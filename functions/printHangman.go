package functions

import (
	"fmt"
	"os"
)

func PrintHangman(attemptLeft int) string {
	w := ""
	if attemptLeft != 10 {
		if Read("hangman.txt")[9] == byte(13) {
			w = (string(Read("hangman.txt")[(9-attemptLeft)*79 : (9-attemptLeft)*79+77]))
		} else {
			w = (string(Read("hangman.txt")[(9-attemptLeft)*71 : (9-attemptLeft)*71+70]))
		}
	}
	return w
}

func Read(s string) []byte {
	file, err := os.Open(s)
	if err != nil {
		fmt.Println(err)
	}
	arr := make([]byte, 1200)
	file.Read(arr)
	return arr
}
