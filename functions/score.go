package functions

// return the score based on the user's attempts left and the difficulty he chose
func Score(remainingAttempts int, difficulty string) int {
	// starting point
	basePoints := 100
	// calculate the number of error made
	errorMade := 10 - remainingAttempts

	// initialize a map for the bonus point
	difficultyPoints := map[string]int{
		"EASY":   10,
		"NORMAL": 30,
		"HARD":   50,
	}

	// calcul of the points
	points := basePoints - (errorMade * 10)

	// adding the bonus
	if bonus, exists := difficultyPoints[difficulty]; exists {
		points += bonus
	}

	if points < 0 {
		points = 0
	}

	return points
}
