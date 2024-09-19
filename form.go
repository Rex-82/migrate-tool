package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"migratetool/utils"
)

func GetCredentialsAndAction(FormData *FormData, theme *huh.Theme) error {

	return huh.NewForm(
		huh.NewGroup(
			// Input for MySQL username
			huh.NewInput().
				Title("MySQL username").
				Value(&formData.username).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("Database name cannot be empty")
					}
					return nil
				}),

			// Input for MySQL password
			huh.NewInput().
				Title("MySQL password").
				Value(&formData.password).EchoMode(huh.EchoModePassword),

			// Input for MySQL db name
			huh.NewInput().
				Title("Database").
				Value(&formData.db).
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
				Value(&formData.action),
		),
	).WithTheme(theme).Run()
}

func GetMigrationType(formData *FormData, theme *huh.Theme) error {

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
				Value(&formData.migrationType),
		),
	).WithTheme(theme).Run()
}

func GetDirectory(formData *FormData, theme *huh.Theme) error {

	var actionTitle string

	switch formData.action {
	case "migrate":
		actionTitle = "Destination directory"
	case "upload":
		actionTitle = "Migrations directory"

	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(actionTitle).
				Value(&formData.directory),
		),
	).WithTheme(theme).Run()

	if err != nil {
		return err
	}

	formData.directory = utils.PathFormat(formData.directory)

	return nil
}

func GetSelectedMigration(formData *FormData, theme *huh.Theme) error {

	sqlFiles, err := utils.ListFiles(formData.directory, ".sql")
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
				Value(&formData.selectedMigration),
		),
	).WithTheme(theme).Run()
}

func DisplaySummary(FormData *FormData, theme *huh.Theme) error {

	titleStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#8D63C9")).Padding(0, 1).Margin(1, 0).Width(32).AlignHorizontal(lipgloss.Center)

	valueStyle := lipgloss.NewStyle().Italic(true).Bold(true).Foreground(lipgloss.Color("#C2B3AC"))

	rowStyle := lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#666666"))
	lastRowStyle := lipgloss.NewStyle()

	diretoryLabel := "Destination:"
	migrationPath := formData.directory

	if formData.action == "upload" {
		diretoryLabel = "Source:"
		formData.selectedMigration = strings.Split(formData.selectedMigration, " ")[0]
		migrationPath = formData.directory + formData.selectedMigration
	}

	// Display the Geted information
	fmt.Println(titleStyle.Render("Migration Information"))
	fmt.Printf("%s %s\n", rowStyle.Render("Username:"), valueStyle.Render(formData.username))
	fmt.Printf("%s %s\n", rowStyle.Render("Password:"), valueStyle.Render(strings.Repeat("*", len(formData.password))))
	fmt.Printf("%s %s\n", rowStyle.Render("Database:"), valueStyle.Render(formData.db))
	fmt.Printf("%s %s\n", rowStyle.Render("Action:"), valueStyle.Render(formData.action))
	fmt.Printf("%s %s\n", rowStyle.Render("Migration Type:"), valueStyle.Render(formData.migrationType))
	fmt.Println(lastRowStyle.Render(rowStyle.Render(diretoryLabel), valueStyle.Render(migrationPath)))

	fmt.Println("\n")

	// Confirmation step before proceeding with the action
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm migration?").
				Affirmative("Confirm").
				Negative("Cancel").
				Value(&formData.confirm),
		),
	).WithTheme(theme).Run()
}
