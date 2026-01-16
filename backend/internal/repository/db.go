package repository

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"backend/internal/model"
)

func openDB(cfg *model.Config) (*sql.DB, error) {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
