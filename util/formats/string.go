package formats

import (
	"encoding/base64"
	"fmt"
	"path"
	"regexp"
	"strings"
)

func ToFileNameString(fileName string) string {
	emojiRegex := regexp.MustCompile(`[\p{So}\p{Sk}\p{Sm}\p{Sc}\p{S}]`)
	fileName = emojiRegex.ReplaceAllString(fileName, "")

	charsToReplace := []string{"/", ".", "*", "?", "\"", "<", ">", "|", "{", "}", "\\", "^", "~", "[", "]", "`", "!", "$", "&", "'", "(", ")", ",", ":", ";", "@"}
	for _, char := range charsToReplace {
		fileName = strings.ReplaceAll(fileName, char, "")
	}

	spaceCharsToReplace := []string{"~", "#", "+", " ", "=", "_"}
	for _, char := range spaceCharsToReplace {
		fileName = strings.ReplaceAll(fileName, char, "-")
	}

	if len(fileName) > 120 {
		fileName = fileName[0:120]
	}

	return strings.ToLower(fileName)
}

func Utf8ToBase64(input string) (string, error) {
	// Convert UTF-8 string to byte slice
	utf8Bytes := []byte(input)

	// Encode byte slice to base64
	base64String := base64.StdEncoding.EncodeToString(utf8Bytes)

	return base64String, nil
}

func ExtractImageNameFromUrl(url string) (imageName, extension string, err error) {
	// Extract the path component of the URL
	imagePath := path.Base(url)

	// Split the filename and extension
	parts := strings.Split(imagePath, ".")

	// Check if the URL has a valid format
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid URL format")
	}

	// Get the image name and extension
	imageName = parts[0]
	extension = parts[1]

	return imageName, extension, nil
}
