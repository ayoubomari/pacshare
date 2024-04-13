package fileDownloader

import (
	"encoding/base64"
	"io"
	"net/http"
)

// this function return a byte slice contain data (file content)
func FetchFileContentAsString(url string) (string, error) {
	// Send an HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(body), nil
}
