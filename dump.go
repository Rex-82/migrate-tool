package main

import (
	"fmt"
	"migratetool/models"
	"migratetool/utils"
	"os/exec"
	"time"
)

func RunMysqldump(command string) error {

	err := utils.EnsureDirExists(models.FormData.Directory)
	if err != nil {
		return err
	}

	// Generate timestamp in the format YYYYMMDD_HHMMSS
	timestamp := time.Now().Format("20060102_150405")

	// Build the file name based on the timestamp and migration type
	var fileName string
	switch models.FormData.MigrationType {
	case "schema":
		fileName = fmt.Sprintf("%s_db_schema.sql", timestamp)
	case "data":
		fileName = fmt.Sprintf("%s_db_data.sql", timestamp)
	case "both":
		fileName = fmt.Sprintf("%s_db_full.sql", timestamp)
	default:
		return fmt.Errorf("invalid migration type: %s", models.FormData.MigrationType)
	}

	// Full path of the dump file
	filePath := fmt.Sprintf("%s%s", models.FormData.Directory, fileName)

	// Construct the mysqldump command based on the migration type
	var dumpCommand string
	switch models.FormData.MigrationType {
	case "schema":
		dumpCommand = fmt.Sprintf("%s -u %s --password='%s' -h %s -P %s --no-data %s > %s", command, models.FormData.Username, models.FormData.Password, models.FormData.Host, models.FormData.Port, models.FormData.Db, filePath)
	case "data":
		dumpCommand = fmt.Sprintf("%s -u %s --password='%s' -h %s -P %s --no-create-info %s > %s", command, models.FormData.Username, models.FormData.Password, models.FormData.Host, models.FormData.Port, models.FormData.Db, filePath)
	case "both":
		dumpCommand = fmt.Sprintf("%s -u %s --password='%s' -h %s -P %s %s > %s", command, models.FormData.Username, models.FormData.Password, models.FormData.Host, models.FormData.Port, models.FormData.Db, filePath)
	}

	fmt.Print("Running...\n")

	cmd := exec.Command(utils.SHELL_CMD, utils.SHELL_CMD_ARG, dumpCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution failed: %v, output: %s", err, output)
	}

	fmt.Printf("Dump file created at: %s\n", filePath)
	return nil
}
