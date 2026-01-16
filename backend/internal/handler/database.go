package handler

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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
		completion_message TEXT,
		review_message TEXT,
		reviewed_by TEXT,
		reviewed_at TEXT,
		id_user INTEGER,
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

	ensureColumns("tasks", map[string]string{
		"completion_message": "TEXT",
		"review_message":     "TEXT",
		"reviewed_by":        "TEXT",
		"reviewed_at":        "TEXT",
	})

	eventTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY,
		title TEXT,
		date TEXT,
		time_range TEXT,
		created_by TEXT,
		description TEXT
	);
	`
	_, err = DB.Exec(eventTable)
	if err != nil {
		log.Fatal(err)
	}

	seedEvents()
}

func ensureColumns(table string, columns map[string]string) {
	for name, columnType := range columns {
		if columnExists(table, name) {
			continue
		}
		query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table, name, columnType)
		if _, err := DB.Exec(query); err != nil {
			log.Fatal(err)
		}
	}
}

func columnExists(table, column string) bool {
	rows, err := DB.Query(fmt.Sprintf("PRAGMA table_info(%s);", table))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var ctype string
		var notNull int
		var defaultValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &ctype, &notNull, &defaultValue, &pk); err != nil {
			log.Fatal(err)
		}
		if name == column {
			return true
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return false
}

func seedEvents() {
	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM events").Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		return
	}

	_, err := DB.Exec(
		`INSERT INTO events (title, date, time_range, created_by, description) VALUES (?, ?, ?, ?, ?)`,
		"Созвон с Марком Олеговичем",
		"18.01.26",
		"17:00-18:30",
		"Слободенюк Никита",
		"",
	)
	if err != nil {
		log.Fatal(err)
	}
}
