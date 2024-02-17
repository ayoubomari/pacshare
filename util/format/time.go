package format

import (
	"strconv"
	"strings"
)

// convert from string:"00:04:50" to Number:290
func DurationStrToSeconds(durationStr string) int {
	if durationStr == "" {
		return 0
	}

	durationArray := strings.Split(durationStr, ":")

	s := 0
	p := 1
	for i := len(durationArray) - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(durationArray[i])
		s += n * p
		p *= 60
	}

	return s
}
