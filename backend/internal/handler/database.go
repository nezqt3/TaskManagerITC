package handler

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"

	"backend/internal/model"
	"backend/internal/logger"
)

var DB *sql.DB

func InitDatabase(cfg *model.Config) {
	logger.Info.Println("Initializing database...")

	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		logger.Fatal.Fatalf("Failed to get absolute path: %v\n", err)
		return
	}
	logger.Info.Printf("Database path resolved: %s\n", absPath)

	DB, err = sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		logger.Fatal.Fatalf("Failed to open database: %v\n", err)
		return
	}
	logger.Info.Println("Database opened successfully")

	createTables()
}

func createTables() {
	logger.Info.Println("Creating tables if not exists...")

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
	if _, err := DB.Exec(userTable); err != nil {
		logger.Fatal.Fatalf("Failed to create 'users' table: %v\n", err)
	} else {
		logger.Info.Println("'users' table ensured")
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
	if _, err := DB.Exec(projectTable); err != nil {
		logger.Fatal.Fatalf("Failed to create 'projects' table: %v\n", err)
	} else {
		logger.Info.Println("'projects' table ensured")
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
	if _, err := DB.Exec(taskTable); err != nil {
		logger.Fatal.Fatalf("Failed to create 'tasks' table: %v\n", err)
	} else {
		logger.Info.Println("'tasks' table ensured")
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
	if _, err := DB.Exec(eventTable); err != nil {
		logger.Fatal.Fatalf("Failed to create 'events' table: %v\n", err)
	} else {
		logger.Info.Println("'events' table ensured")
	}

	seedEvents()
}

func ensureColumns(table string, columns map[string]string) {
	for name, columnType := range columns {
		if columnExists(table, name) {
			logger.Info.Printf("Column '%s' already exists in table '%s'\n", name, table)
			continue
		}
		query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table, name, columnType)
		if _, err := DB.Exec(query); err != nil {
			logger.Fatal.Fatalf("Failed to add column '%s' to table '%s': %v\n", name, table, err)
		} else {
			logger.Info.Printf("Added column '%s' to table '%s'\n", name, table)
		}
	}
}

func columnExists(table, column string) bool {
	rows, err := DB.Query(fmt.Sprintf("PRAGMA table_info(%s);", table))
	if err != nil {
		logger.Fatal.Fatalf("Failed to query table info for '%s': %v\n", table, err)
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
			logger.Fatal.Fatalf("Failed to scan table info for '%s': %v\n", table, err)
		}
		if name == column {
			return true
		}
	}

	if err := rows.Err(); err != nil {
		logger.Fatal.Fatalf("Error iterating table info for '%s': %v\n", table, err)
	}

	return false
}

func seedEvents() {
	logger.Info.Println("Seeding events table if empty...")

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM events").Scan(&count); err != nil {
		logger.Fatal.Fatalf("Failed to count events: %v\n", err)
	}
	logger.Info.Printf("Events table contains %d rows\n", count)

	if count > 0 {
		logger.Info.Println("Events table already seeded")
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
		logger.Fatal.Fatalf("Failed to seed events: %v\n", err)
	} else {
		logger.Info.Println("Events table seeded successfully")
	}
}
