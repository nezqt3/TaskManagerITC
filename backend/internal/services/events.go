package services

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"backend/internal/logger"
	"backend/internal/model"
)

func GetEvents(cfg *model.Config) ([]model.Event, error) {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		logger.Error.Printf("GetEvents: failed to get absolute path for database '%s': %v\n", cfg.NAME_OF_DATABASE, err)
		return nil, err
	}

	db, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		logger.Error.Printf("GetEvents: failed to open database '%s': %v\n", absPath, err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT id, title, date, time_range, created_by, COALESCE(description, '')
		FROM events
		ORDER BY id DESC
	`)
	if err != nil {
		logger.Error.Printf("GetEvents: failed to query events: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	events := make([]model.Event, 0)
	for rows.Next() {
		var e model.Event
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Date,
			&e.TimeRange,
			&e.CreatedBy,
			&e.Description,
		); err != nil {
			logger.Error.Printf("GetEvents: failed to scan row: %v\n", err)
			return nil, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		logger.Error.Printf("GetEvents: rows iteration error: %v\n", err)
		return nil, err
	}

	return events, nil
}
