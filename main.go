package main

import (
	"fmt"
	"log"
	"migratetool/utils"
	"os/exec"

	"github.com/charmbracelet/huh"
)

type FormData struct {
	username          string
	password          string
	db                string
	action            string
	migrationType     string
	selectedMigration string
	confirm           bool
	directory         string
}

var formData = FormData{
	username:  "root",
	directory: "./db/migrations/",
}

var theme = huh.ThemeDracula()

func main() {
	err := GetCredentialsAndAction(&formData, theme)
	if err != nil {
		log.Fatal(err)
	}

	if formData.action == "migrate" {
		err = GetMigrationType(&formData, theme)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = GetDirectory(&formData, theme)
	if err != nil {
		log.Fatal(err)
	}

	if formData.action == "upload" {

		err = GetSelectedMigration(&formData, theme)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = DisplaySummary(&formData, theme)
	if err != nil {
		log.Fatal(err)
	}

	if !formData.confirm {
		fmt.Println("Action cancelled by the user.")
		return
	}

	if formData.action == "migrate" {

		cmd, err := exec.LookPath("mysqldump")
		if err != nil {
			cmd, err = exec.LookPath(utils.MYSQLDUMP_BIN)
			if err != nil {
				log.Fatalf("Command 'mysqldump' not found at: %v", utils.MYSQLDUMP_BIN)
			}
			fmt.Printf(cmd)

		}

		err = RunMysqldump(cmd, formData.username, formData.password, formData.db, formData.migrationType, formData.directory)
		if err != nil {
			log.Fatalf("Failed to run mysqldump: %v", err)
		}
	} else if formData.action == "upload" {

		migrationFile := formData.directory + formData.selectedMigration

		cmd, err := exec.LookPath("mysql")
		if err != nil {
			cmd, err = exec.LookPath(utils.MYSQL_BIN)
			if err != nil {
				log.Fatalf("Command 'mysql' not found at: %v", utils.MYSQL_BIN)
			}
			fmt.Printf(cmd)

		}

		err = applyMigration(cmd, migrationFile)
		if err != nil {
			log.Fatal(err)
		}
	}

}
