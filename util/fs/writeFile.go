package fs

import (
	"os"
)

func WriteFile(filename string, content *string) error {
	// Write the content to the file
	err := os.WriteFile(filename, []byte(*content), 0644)
	if err != nil {
		return err
	}
	return nil
}
