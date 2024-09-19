package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func ListFiles(dir string, extension string) ([]string, error) {

	var sqlFiles []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	// Variables to track the latest file
	var latestFile string
	var latestModTime time.Time

	// Iterate over files to find .sql files and determine the latest one
	for _, file := range files {
		if filepath.Ext(file.Name()) == extension {
			fileInfo, err := file.Info()
			if err != nil {
				return nil, fmt.Errorf("error getting file info: %v", err)
			}

			// Check if this file is the latest one based on modification time
			if latestFile == "" || fileInfo.ModTime().After(latestModTime) {
				latestFile = file.Name()
				latestModTime = fileInfo.ModTime()
			}

			// Add file to the list (will update with "(latest)" later if needed)
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	// Mark the latest file with "(latest)"
	for i, file := range sqlFiles {
		if file == latestFile {
			sqlFiles[i] += " (latest)"
			break
		}
	}

	if len(sqlFiles) == 0 {
		return nil, fmt.Errorf("no SQL files found in the directory: %s", dir)
	}

	return sqlFiles, nil
}
