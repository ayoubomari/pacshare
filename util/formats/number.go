package formats

import (
	"errors"
	"regexp"
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

func ParseFloat(input string) (float64, error) {
	// Define the regular expression pattern
	regexPattern := `([0-9.]+)`

	// Compile the regular expression
	re := regexp.MustCompile(regexPattern)

	// Find the match in the input string
	match := re.FindStringSubmatch(input)

	if len(match) < 2 {
		return 0, errors.New("no number found in the input string")
	}

	// Convert the matched string to a float64
	number, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}
