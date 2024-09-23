package main

import (
	"fmt"
	"migratetool/utils"
	"os/exec"
	"time"
)

func RunMysqldump(command, username, password, databaseName, migrationType, directory string) error {

	err := utils.EnsureDirExists(directory)
	if err != nil {
		return err
	}

	// Generate timestamp in the format YYYYMMDD_HHMMSS
	timestamp := time.Now().Format("20060102_150405")

	// Build the file name based on the timestamp and migration type
	var fileName string
	switch migrationType {
	case "schema":
		fileName = fmt.Sprintf("%s_db_schema.sql", timestamp)
	case "data":
		fileName = fmt.Sprintf("%s_db_data.sql", timestamp)
	case "both":
		fileName = fmt.Sprintf("%s_db_full.sql", timestamp)
	default:
		return fmt.Errorf("invalid migration type: %s", migrationType)
	}

	// Full path of the dump file
	filePath := fmt.Sprintf("%s/%s", directory, fileName)

	// Construct the mysqldump command based on the migration type
	var dumpCommand string
	switch migrationType {
	case "schema":
		dumpCommand = fmt.Sprintf("%s -u %s --password='%s' --no-data %s > %s", command, username, password, databaseName, filePath)
	case "data":
		dumpCommand = fmt.Sprintf("%s -u %s --password='%s' --no-create-info %s > %s", command, username, password, databaseName, filePath)
	case "both":
		dumpCommand = fmt.Sprintf("%s -u %s --password='%s' %s > %s", command, username, password, databaseName, filePath)
	}

	fmt.Print("Running...\n")

	cmd := exec.Command("sh", "-c", dumpCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution failed: %v, output: %s", err, output)
	}

	fmt.Printf("Dump file created at: %s\n", filePath)
	return nil
}
