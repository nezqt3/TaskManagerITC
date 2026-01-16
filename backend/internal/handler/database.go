package handler

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"path/filepath"

	"backend/internal/model"
)

var DB *sql.DB

func InitDatabase(cfg *model.Config) {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)

	if err != nil {
		log.Fatal(err)
		return
	}

	DB, err = sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	createTables()
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		TelegramID     TEXT PRIMARY KEY,
		FirstName      TEXT,
		LastName       TEXT, 
		Username       TEXT, 
		PhotoURL       TEXT, 
		FullName       TEXT, 
		DateOfBirthday TEXT, 
		NumberOfPhone  TEXT, 
		Role           TEXT, 
		MayToOpen      BOOLEAN   
	);
	`

	_, err := DB.Exec(userTable)
	if err != nil {
		log.Fatal(err)
	}

	projectTable := `
	CREATE TABLE IF NOT EXISTS projects (
		id     			INTEGER PRIMARY KEY,
		description     TEXT,
		users       	TEXT, 
		title       	TEXT, 
		status      	TEXT  
	);
	`

	_, err = DB.Exec(projectTable)
	if err != nil {
		log.Fatal(err)
	}

	taskTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id    INTEGER PRIMARY KEY, 
		description TEXT, 
		deadline DATE, 
		status TEXT, 
        user_id INTEGER,
		user TEXT, 
		title TEXT,
		author TEXT,
		id_project INTEGER
	);
	`
	_, err = DB.Exec(taskTable)
	if err != nil {
		log.Fatal(err)
	}
}