package main

import (
	"fmt"
	"log"
	"migratetool/models"
	"migratetool/utils"
	"os/exec"

	"github.com/charmbracelet/huh"
)

var theme = huh.ThemeDracula()

func main() {
	err := GetServerInfo(&models.FormData, theme)
	if err != nil {
		log.Fatal(err)
	}

	err = GetCredentialsAndAction(&models.FormData, theme)
	if err != nil {
		log.Fatal(err)
	}

	if models.FormData.Action == "migrate" {
		err = GetMigrationType(&models.FormData, theme)
		if err != nil {
			log.Fatal(err)
		}

		switch models.FormData.MigrationType {
		case "schema":
			models.FormData.Directory += "schema/"
		case "data":
			models.FormData.Directory += "data/"
		}

	}

	err = GetDirectory(models.FormData, theme)
	if err != nil {
		log.Fatal(err)
	}

	if models.FormData.Action == "upload" {

		err = GetSelectedMigration(models.FormData, theme)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = DisplaySummary(models.FormData, theme)
	if err != nil {
		log.Fatal(err)
	}

	if !models.FormData.Confirm {
		fmt.Println("Action cancelled by the user.")
		return
	}

	if models.FormData.Action == "migrate" {

		cmd, err := exec.LookPath("mysqldump")
		if err != nil {
			cmd, err = exec.LookPath(utils.MYSQLDUMP_BIN)
			if err != nil {
				log.Fatalf("Command 'mysqldump' not found at: %v", utils.MYSQLDUMP_BIN)
			}
			fmt.Printf(cmd)

		}

		err = RunMysqldump(cmd)
		if err != nil {
			log.Fatalf("Failed to run mysqldump: %v", err)
		}
	} else if models.FormData.Action == "upload" {

		migrationFile := models.FormData.Directory + models.FormData.SelectedMigration

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
