package fileDownloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ayoubomari/pacshare/util/fs"
)

// download the entire file
func DownloadEtireFile(mediaUrl string, filePath string) error {
	// Create an output file (you can use any io.Writer)
	outputFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("DownloadEtireFile: %w", err)
	}
	defer outputFile.Close()

	// Initialize HTTP client
	client := &http.Client{}

	// Create HTTP GET request
	req, err := http.NewRequest("GET", mediaUrl, nil)
	if err != nil {
		fs.DeleteFile(filePath)
		return fmt.Errorf("DownloadEtireFile: %w", err)
	}

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fs.DeleteFile(filePath)
		return fmt.Errorf("DownloadEtireFile: %w", err)
	}
	defer resp.Body.Close()

	// if file is (apk, obb) add one byte to the first of the file to make it elegible on facebook
	if strings.Contains(filePath, "_pac.apk") || strings.Contains(filePath, "_pac.obb") {
		_, err := outputFile.Write([]byte("0"))
		if err != nil {
			fs.DeleteFile(filePath)
			return fmt.Errorf("DownloadEtireFile: %w", err)
		}
	}

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		fs.DeleteFile(filePath)
		return fmt.Errorf("DownloadEtireFile: %w", err)
	}

	return nil
}
