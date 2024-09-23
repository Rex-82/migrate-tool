package main

import (
	"fmt"
	"os"
	"os/exec"
)

func applyMigration(command string, dumpFile string) error {
	fmt.Println("Applying migration...")

	var cmd *exec.Cmd

	file, err := os.Open(dumpFile)
	if err != nil {
		return fmt.Errorf("failed to open dump file: %v", err)
	}

	defer file.Close()

	cmd = exec.Command(command, "-u", formData.username, "--password="+formData.password, formData.db)
	cmd.Stdin = file // Redirect the file content to MySQL's stdin

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to apply migration: %v", err)
	}

	fmt.Println("Migration successfully applied to the database.")
	return nil
}
