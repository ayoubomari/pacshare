package fs

import (
	"fmt"
	"os"
)

func DeleteFile(filePath string) error {
	// Attempt to remove the file
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
