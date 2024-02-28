package fileDownloader

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/ayoubomari/pacshare/util/fs"
	"github.com/gofiber/fiber/v2"
)

// download the entire file
func downloadEtireFileWithFiber(mediaUrl string, filePath string) error {
	outputFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("DownloadEtireFileWithFiber: %w", err)
	}
	defer outputFile.Close()

	agent := fiber.Get(mediaUrl)
	_, body, errs := agent.Bytes()
	if len(errs) > 0 {
		fs.DeleteFile(filePath)
		return errs[0]
	}

	bodyReader := bytes.NewReader(body)
	_, err = io.Copy(outputFile, bodyReader)
	if err != nil {
		fs.DeleteFile(filePath)
		return fmt.Errorf("DownloadEtireFileWithFiber: %w", err)
	}

	return nil
}
