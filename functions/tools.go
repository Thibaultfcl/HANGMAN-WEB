package functions

// function to check if a specific int is on a slice of int
func Contains(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// function to check if a specific string is on a slice of string
func containsStr(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}
