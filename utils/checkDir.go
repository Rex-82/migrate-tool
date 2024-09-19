package utils

import (
	"fmt"
	"os"
)

func EnsureDirExists(dir string) error {
	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory does not exist, create it
		fmt.Printf("Directory '%s' not found. Creating it...\n", dir)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	} else if err != nil {
		// An error other than "not exists" occurred
		return fmt.Errorf("error checking directory: %v", err)
	}

	return nil
}
