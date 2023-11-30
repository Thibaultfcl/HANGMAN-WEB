package functions

// func GetUserInput() string {
// 	var input string
// 	fmt.Scan(&input)
// 	return input
// }

func Contains(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// func containsStr(slice []string, item string) bool {
// 	for _, i := range slice {
// 		if i == item {
// 			return true
// 		}
// 	}
// 	return false
// }
