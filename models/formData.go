package models

type FormDataType struct {
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
	Username:  "root",
	Directory: "./db/migrations/",
}
