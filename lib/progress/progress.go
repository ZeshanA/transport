package progress

import (
	"fmt"
	"math"
)

// PrintAtIntervals prints a progress message every time another ~10%
// of inputs are processed
func PrintAtIntervals(completed int, total int, process string) {
	tenPercentIncrement := int(math.Floor(float64(total / 10)))
	if tenPercentIncrement == 0 || completed%tenPercentIncrement == 0 {
		fmt.Printf(
			"%s: processed %d/%d tasks (~%f%%)\n",
			process, completed, total, float64(completed)/float64(total)*100,
		)
	}
}
