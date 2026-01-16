package repository

import "backend/internal/model"

type ProjectRow struct {
	ID          int
	Title       string
	Description string
	Status      string
	RawUsers    string
}

func GetProjects(cfg *model.Config) ([]ProjectRow, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, description, status, users FROM projects ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]ProjectRow, 0)
	for rows.Next() {
		var row ProjectRow
		if err := rows.Scan(&row.ID, &row.Title, &row.Description, &row.Status, &row.RawUsers); err != nil {
			return nil, err
		}
		projects = append(projects, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func UpdateProjectMembers(cfg *model.Config, projectID int, usersJSON string) error {
	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`UPDATE projects SET users = ? WHERE id = ?`, usersJSON, projectID)
	return err
}

func UpdateProjectStatus(cfg *model.Config, projectID int, status string) error {
	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`UPDATE projects SET status = ? WHERE id = ?`, status, projectID)
	return err
}

func GetTaskStatusesByProjectID(cfg *model.Config, projectID int) ([]string, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT status FROM tasks WHERE id_project = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := []string{}
	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statuses, nil
}
