package handler

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"path/filepath"
	"log"

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
}

func getAllUsers() []model.UserProfile {
    rows, err := DB.Query("SELECT * FROM users")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var users []model.UserProfile
    for rows.Next() {
        var u model.UserProfile
        err := rows.Scan(
            &u.TelegramID,
            &u.FirstName,
            &u.LastName,
            &u.Username,
            &u.PhotoURL,
            &u.FullName,
            &u.DateOfBirthday,
            &u.NumberOfPhone,
            &u.Role,
            &u.MayToOpen,
        )
        if err != nil {
            log.Fatal(err)
        }
        users = append(users, u)
    }
    return users
}

func createUser(user *model.UserProfile) {
    _, err := DB.Exec(
        "INSERT INTO users VALUES(?,?,?,?,?,?,?,?,?,?)",
        user.TelegramID,
        user.FirstName,
        user.LastName,
        user.Username,
        user.PhotoURL,
        user.FullName,
        user.DateOfBirthday,
        user.NumberOfPhone,
        user.Role,
        user.MayToOpen,
    )
    if err != nil {
        log.Fatal(err)
    }
}

func getUserByTelegramID(telegramID string) (*model.UserProfile, error) {
    u := &model.UserProfile{}
    err := DB.QueryRow("SELECT * FROM users WHERE TelegramID = ?", telegramID).Scan(
        &u.TelegramID,
        &u.FirstName,
        &u.LastName,
        &u.Username,
        &u.PhotoURL,
        &u.FullName,
        &u.DateOfBirthday,
        &u.NumberOfPhone,
        &u.Role,
        &u.MayToOpen,
    )
    if err != nil {
        return nil, err
    }
    return u, nil
}