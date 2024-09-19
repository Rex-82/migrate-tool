package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	username          string = "root"
	password          string
	db                string
	action            string
	migrationType     string
	selectedMigration string
	confirm           bool
	directory         string = "./db/migrations"
)

func main() {
	// First step: Collect MySQL credentials and action
	err := huh.NewForm(
		huh.NewGroup(
			// Input for MySQL username
			huh.NewInput().
				Title("MySQL username").
				Value(&username),

			// Input for MySQL password
			huh.NewInput().
				Title("MySQL password").
				Value(&password).EchoMode(huh.EchoModePassword),

			// Input for MySQL db name
			huh.NewInput().
				Title("Database").
				Value(&db).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("Database name cannot be empty")
					}
					return nil
				}),
		),
		huh.NewGroup(
			// Selection for Action
			huh.NewSelect[string]().
				Title("Action").
				Options(
					huh.NewOption("migrate", "migrate"),
					huh.NewOption("upload", "upload"),
				).
				Value(&action),
		),
	).Run()
	if err != nil {
		log.Fatal(err)
	}

	if action == "migrate" {
		// Second step: Collect migration type
		err = huh.NewForm(
			huh.NewGroup(
				// Selection for Migration Type
				huh.NewSelect[string]().
					Title("Migration Type").
					Options(
						huh.NewOption("schema", "schema"),
						huh.NewOption("data", "data"),
						huh.NewOption("both", "both"),
					).
					Value(&migrationType),
			),
		).Run()
		if err != nil {
			log.Fatal(err)
		}

		// Third step: Collect destination directory if action is "migrate"
		err = huh.NewForm(
			huh.NewGroup(
				// Input for Destination Directory
				huh.NewInput().
					Title("Destination directory").
					Value(&directory),
			),
		).Run()
		if err != nil {
			log.Fatal(err)
		}
	} else if action == "upload" {
		// Step for upload action: Collect directory and list available migrations
		err = huh.NewForm(
			huh.NewGroup(
				// Input for the directory where migrations are stored
				huh.NewInput().
					Title("Migrations directory").
					Value(&directory),
			),
		).Run()
		if err != nil {
			log.Fatal(err)
		}

		// List all .sql migration files in the provided directory
		sqlFiles, err := listSQLFiles(directory)
		if err != nil {
			log.Fatal(err)
		}

		// Prompt user to select a migration file from the available list
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select a migration file").
					Options(
						func() (opts []huh.Option[string]) {
							for _, file := range sqlFiles {
								opts = append(opts, huh.NewOption(file, file))
							}
							return
						}()...,
					).
					Value(&selectedMigration),
			),
		).Run()
		if err != nil {
			log.Fatal(err)
		}

	}

	titleStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#6666BB")).Padding(0, 1).Margin(1, 0)

	valueStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#CCCCDD"))

	rowStyle := lipgloss.NewStyle().PaddingLeft(1)
	lastRowStyle := lipgloss.NewStyle().MarginBottom(1)

	diretoryLabel := "Destination:"
	if action == "upload" {
		diretoryLabel = "Source:"
		directory += strings.Split(selectedMigration, " ")[0]
	}

	// Display the collected information
	fmt.Println(titleStyle.Render("Migration Information:"))
	fmt.Printf("%s %s\n", rowStyle.Render("Username:"), valueStyle.Render(username))
	fmt.Printf("%s %s\n", rowStyle.Render("Password:"), valueStyle.Render(strings.Repeat("*", len(password))))
	fmt.Printf("%s %s\n", rowStyle.Render("Database:"), valueStyle.Render(db))
	fmt.Printf("%s %s\n", rowStyle.Render("Action:"), valueStyle.Render(action))
	fmt.Printf("%s %s\n", rowStyle.Render("Migration Type:"), valueStyle.Render(migrationType))
	fmt.Println(lastRowStyle.Render(rowStyle.Render(diretoryLabel), valueStyle.Render(directory)))

	fmt.Println("\n")

	// Confirmation step before proceeding with the action
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm migration?").
				Affirmative("Confirm!").
				Negative("Cancel").
				Value(&confirm),
		),
	).Run()

	if err != nil {
		log.Fatal(err)
	}

	if !confirm {
		fmt.Println("Action cancelled by the user.")
		return
	}

	// Construct and run the mysqldump command
	if action == "migrate" {
		err = runMysqldump(username, password, db, migrationType, directory)
		if err != nil {
			log.Fatalf("Failed to run mysqldump: %v", err)
		}
	}

	// Apply the selected migration file
	if action == "upload" {
		err = applyMigration(filepath.Join(directory, selectedMigration))
		if err != nil {
			log.Fatal(err)
		}
	}

}

// Generic function to ensure a directory exists, creating it if necessary
func ensureDirExists(dir string) error {
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

// runMysqldump constructs and executes the mysqldump command
func runMysqldump(username, password, databaseName, migrationType, directory string) error {

	// Ensure the "db/migrations" directory exists
	if err := ensureDirExists(directory); err != nil {
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
		dumpCommand = fmt.Sprintf("mysqldump -u %s -password='%s' --no-data %s > %s", username, password, databaseName, filePath)
	case "data":
		dumpCommand = fmt.Sprintf("mysqldump -u %s -password='%s' --no-create-info %s > %s", username, password, databaseName, filePath)
	case "both":
		dumpCommand = fmt.Sprintf("mysqldump -u %s -password='%s' %s > %s", username, password, databaseName, filePath)
	}

	fmt.Printf("Running command: %s\n", dumpCommand)

	cmd := exec.Command("sh", "-c", dumpCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution failed: %v, output: %s", err, output)
	}

	fmt.Printf("Dump file created at: %s\n", filePath)
	return nil
}

// Function to apply the migration to the database
func applyMigration(dumpFile string) error {
	fmt.Println("Applying migration...")

	// Run the mysql command to import the SQL dump file back into the database
	cmd := exec.Command("mysql", "-u", username, "-p"+password, db, "-e", fmt.Sprintf("source %s", dumpFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to apply migration: %v", err)
	}

	fmt.Println("Migration successfully applied to the database.")
	return nil
}

// Function to list all .sql files in a given directory
func listSQLFiles(dir string) ([]string, error) {
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
		if filepath.Ext(file.Name()) == ".sql" {
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
