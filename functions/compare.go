package functions

//compare the letter entered with the actual game data
func Compare(HiddenWord []string, AttemptLeft int, WordToGuess string, letter string, letterUsed []string) ([]string, int, []string) {
	if len(letter) == 1 { // letter
		//the letter has already been used
		if containsStr(letterUsed, letter) {
			return HiddenWord, AttemptLeft, letterUsed
		}

		// here we check if the letter is on the secret word and if it hasn't been already found
		found := false
		for i, char := range WordToGuess {
			if char == rune(letter[0]) && HiddenWord[i] == "_" {
				// if it the case we remplace the hidden word by the letter
				HiddenWord[i] = string(char)
				found = true
			}
		}

		if !found {
			// if it hasn't been found the attempts left decreases by one
			AttemptLeft--
		}
	} else if letter == WordToGuess { // word
		for i, char := range WordToGuess {
			//the hidden word is set to the secret one
			HiddenWord[i] = string(char)
		}
	} else {
		// if it hasn't been found the attempts left decreases by two
		AttemptLeft--
		AttemptLeft--
	}

	//add the letter
	letterUsed = append(letterUsed, letter)
	//send back the data needed for the verification
	return HiddenWord, AttemptLeft, letterUsed
}
