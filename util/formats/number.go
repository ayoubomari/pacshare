package formats

import (
	"strconv"
	"strings"
)

func FormatNumberWithCommas(nmbr int64) string {
	// Convert nmbr to a string
	nmbrStr := strconv.FormatInt(nmbr, 10)

	// Use strings.Builder to efficiently build the result
	var result strings.Builder

	// Iterate over each character in the string
	for i, char := range nmbrStr {
		// Add a comma after every group of three digits, except for the last group
		if i > 0 && (len(nmbrStr)-i)%3 == 0 {
			result.WriteRune(',')
		}

		// Append the current character to the result
		result.WriteRune(char)
	}

	// Return the formatted string
	return result.String()
}
