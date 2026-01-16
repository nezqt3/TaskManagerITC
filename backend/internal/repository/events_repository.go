package repository

import "backend/internal/model"

func GetEvents(cfg *model.Config) ([]model.Event, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT id, title, date, time_range, created_by, COALESCE(description, '')
		FROM events
		ORDER BY id DESC
	`)
	if err != nil {
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
			return nil, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
