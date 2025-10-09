package util

import (
	"fmt"
	"io"
	"os"
)

// Writes content to file at specified file path
func WriteToFile(content []string, path string) error {

	// File gets truncated if it already exists
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var count int = 0
	for _, entry := range content {
		f.Write([]byte((entry)))
		count++
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	fmt.Println("Saved file to", path)
	return nil
}
