package service

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"backend/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

func GetUsers(cfg *model.Config) ([]model.UserProfile, error) {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT TelegramID, FirstName, LastName, Username, PhotoURL, FullName, DateOfBirthday, NumberOfPhone, Role, MayToOpen
		FROM users
		ORDER BY FullName
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]model.UserProfile, 0)
	for rows.Next() {
		var u model.UserProfile
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByTelegramID(cfg *model.Config, telegramID string) (*model.UserProfile, error) {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	normalizedID := strings.TrimSpace(telegramID)
	u := &model.UserProfile{}
	err = db.QueryRow(`
		SELECT TelegramID, FirstName, LastName, Username, PhotoURL, FullName, DateOfBirthday, NumberOfPhone, Role, MayToOpen
		FROM users
		WHERE TelegramID = ?
	`, normalizedID).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}
