package formats

import (
	"regexp"
	"strings"
)

func ToFileNameString(fileName string) string {
	emojiRegex := regexp.MustCompile(`[\p{So}\p{Sk}\p{Sm}\p{Sc}\p{S}]`)
	fileName = emojiRegex.ReplaceAllString(fileName, "")

	charsToReplace := []string{"/", "*", "?", "\"", "<", ">", "|", "{", "}", "\\", "^", "~", "[", "]", "`", "!", "$", "&", "'", "(", ")", ",", ":", ";", "@"}
	for _, char := range charsToReplace {
		fileName = strings.ReplaceAll(fileName, char, "")
	}

	spaceCharsToReplace := []string{"~", "#", "+", " ", "=", "_"}
	for _, char := range spaceCharsToReplace {
		fileName = strings.ReplaceAll(fileName, char, "-")
	}
	return strings.ToLower(fileName)
}
