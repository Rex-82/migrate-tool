package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"migratetool/models"
	"migratetool/utils"
)

func GetServerInfo(FormData *models.FormDataType, theme *huh.Theme) error {

	return huh.NewForm(
		huh.NewGroup(
			// Input for server host
			huh.NewInput().
				Title("Host").
				Value(&models.FormData.Host).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("Host name cannot be empty")
					}
					return nil
				}),

			// Input for host port
			huh.NewInput().
				Title("Port").
				Value(&models.FormData.Port).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("Port cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()
}

func GetCredentialsAndAction(FormData *models.FormDataType, theme *huh.Theme) error {

	return huh.NewForm(
		huh.NewGroup(
			// Input for MySQL username
			huh.NewInput().
				Title("MySQL username").
				Value(&models.FormData.Username).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("Database name cannot be empty")
					}
					return nil
				}),

			// Input for MySQL password
			huh.NewInput().
				Title("MySQL password").
				Value(&models.FormData.Password).EchoMode(huh.EchoModePassword),

			// Input for MySQL db name
			huh.NewInput().
				Title("Database").
				Value(&models.FormData.Db).
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
				Value(&models.FormData.Action),
		),
	).WithTheme(theme).Run()
}

func GetMigrationType(FormData *models.FormDataType, theme *huh.Theme) error {

	return huh.NewForm(
		huh.NewGroup(
			// Selection for Migration Type
			huh.NewSelect[string]().
				Title("Migration Type").
				Options(
					huh.NewOption("schema", "schema"),
					huh.NewOption("data", "data"),
					huh.NewOption("both", "both"),
				).
				Value(&models.FormData.MigrationType),
		),
	).WithTheme(theme).Run()
}

func GetDirectory(FormData models.FormDataType, theme *huh.Theme) error {

	var actionTitle string

	switch models.FormData.Action {
	case "migrate":
		actionTitle = "Destination directory"
	case "upload":
		actionTitle = "Migrations directory"

	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(actionTitle).
				Value(&models.FormData.Directory),
		),
	).WithTheme(theme).Run()

	if err != nil {
		return err
	}

	models.FormData.Directory = utils.PathFormat(models.FormData.Directory)

	return nil
}

func GetSelectedMigration(FormData models.FormDataType, theme *huh.Theme) error {

	sqlFiles, err := utils.ListFiles(models.FormData.Directory, ".sql")
	if err != nil {
		log.Fatal(err)
	}

	// Prompt user to select a migration file from the available list
	return huh.NewForm(
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
				Value(&models.FormData.SelectedMigration),
		),
	).WithTheme(theme).Run()
}

func DisplaySummary(FormData models.FormDataType, theme *huh.Theme) error {

	titleStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#8D63C9")).Padding(0, 1).Margin(1, 0).Width(32).AlignHorizontal(lipgloss.Center)

	valueStyle := lipgloss.NewStyle().Italic(true).Bold(true).Foreground(lipgloss.Color("#C2B3AC"))

	rowStyle := lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#666666"))
	lastRowStyle := lipgloss.NewStyle()

	diretoryLabel := "Destination:"
	migrationPath := models.FormData.Directory

	if models.FormData.Action == "upload" {
		diretoryLabel = "Source:"
		models.FormData.SelectedMigration = strings.Split(models.FormData.SelectedMigration, " ")[0]
		migrationPath = models.FormData.Directory + models.FormData.SelectedMigration
	}

	// Display the Geted information
	fmt.Println(titleStyle.Render("Migration Information"))
	fmt.Printf("%s %s\n", rowStyle.Render("Host:"), valueStyle.Render(models.FormData.Host))
	fmt.Printf("%s %s\n", rowStyle.Render("Port:"), valueStyle.Render(models.FormData.Port))
	fmt.Printf("%s %s\n", rowStyle.Render("Username:"), valueStyle.Render(models.FormData.Username))
	fmt.Printf("%s %s\n", rowStyle.Render("Password:"), valueStyle.Render(strings.Repeat("*", len(models.FormData.Password))))
	fmt.Printf("%s %s\n", rowStyle.Render("Database:"), valueStyle.Render(models.FormData.Db))
	fmt.Printf("%s %s\n", rowStyle.Render("Action:"), valueStyle.Render(models.FormData.Action))
	fmt.Printf("%s %s\n", rowStyle.Render("Migration Type:"), valueStyle.Render(models.FormData.MigrationType))
	fmt.Println(lastRowStyle.Render(rowStyle.Render(diretoryLabel), valueStyle.Render(migrationPath)))

	fmt.Println("\n")

	// Confirmation step before proceeding with the action
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm migration?").
				Affirmative("Confirm").
				Negative("Cancel").
				Value(&models.FormData.Confirm),
		),
	).WithTheme(theme).Run()
}
