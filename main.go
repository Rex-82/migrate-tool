package main

import (
	"fmt"
	"log"

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
		err = RunMysqldump(formData.username, formData.password, formData.db, formData.migrationType, formData.directory)
		if err != nil {
			log.Fatalf("Failed to run mysqldump: %v", err)
		}
	} else if formData.action == "upload" {

		migrationFile := formData.directory + formData.selectedMigration

		err = applyMigration(migrationFile)
		if err != nil {
			log.Fatal(err)
		}
	}

}
