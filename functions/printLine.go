package functions

import (
	"bufio"
	"fmt"
	"os"
)

func PrintLine(file *os.File, lineCounter int, linesToDisplay int) {
	scanner := bufio.NewScanner(file)
	linesDisplayed := 0

	for scanner.Scan() {
		if linesDisplayed < lineCounter {
			linesDisplayed++
			continue
		}

		if linesDisplayed >= lineCounter+linesToDisplay {
			break
		}

		fmt.Printf("%s\n", scanner.Text())
		linesDisplayed++
	}
}
