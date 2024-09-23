package models

type FormDataType struct {
	Host              string
	Port              string
	Username          string
	Password          string
	Db                string
	Action            string
	MigrationType     string
	SelectedMigration string
	Confirm           bool
	Directory         string
}

var FormData = FormDataType{
	Host:      "localhost",
	Port:      "3306",
	Username:  "root",
	Directory: "./db/migrations/",
}
