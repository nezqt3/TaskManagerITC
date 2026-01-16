package repository

import (
	"database/sql"
	"time"

	"backend/internal/model"
)

func CreateTask(cfg *model.Config, task *model.Task) (int, error) {
	db, err := openDB(cfg)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	result, err := db.Exec(`
		INSERT INTO tasks (
			description,
			deadline,
			status,
			completion_message,
			review_message,
			reviewed_by,
			reviewed_at,
			id_user,
			user,
			title,
			author,
			id_project
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.Description,
		task.Deadline,
		task.Status,
		"",
		"",
		"",
		"",
		task.IdUser,
		task.User,
		task.Title,
		task.Author,
		task.IdProject,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetTasksByProjectID(cfg *model.Config, projectID int) ([]model.Task, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT id, description, deadline, status,
			COALESCE(completion_message, ''),
			COALESCE(review_message, ''),
			COALESCE(reviewed_by, ''),
			COALESCE(reviewed_at, ''),
			user, title, author, id_project, COALESCE(id_user, 0)
		FROM tasks
		WHERE id_project = ?
		ORDER BY id DESC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		err := rows.Scan(
			&t.ID,
			&t.Description,
			&t.Deadline,
			&t.Status,
			&t.CompletionMessage,
			&t.ReviewMessage,
			&t.ReviewedBy,
			&t.ReviewedAt,
			&t.User,
			&t.Title,
			&t.Author,
			&t.IdProject,
			&t.IdUser,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func UpdateTask(cfg *model.Config, task *model.Task) error {
	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE tasks
		SET title = ?, description = ?, deadline = ?, status = ?, user = ?, id_user = ?
		WHERE id = ?
	`,
		task.Title,
		task.Description,
		task.Deadline,
		task.Status,
		task.User,
		task.IdUser,
		task.ID,
	)
	return err
}

func DeleteTask(cfg *model.Config, taskID int) error {
	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM tasks WHERE id = ?`, taskID)
	return err
}

func SubmitTaskCompletion(cfg *model.Config, taskID int, message string) error {
	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE tasks
		SET status = ?, completion_message = ?, review_message = '', reviewed_by = '', reviewed_at = ''
		WHERE id = ?
	`, "На проверке", message, taskID)
	return err
}

func ReviewTaskCompletion(cfg *model.Config, taskID int, approved bool, reviewer string, message string) error {
	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	status := "Отклонена"
	if approved {
		status = "Выполнена"
	}

	_, err = db.Exec(`
		UPDATE tasks
		SET status = ?, review_message = ?, reviewed_by = ?, reviewed_at = ?
		WHERE id = ?
	`, status, message, reviewer, time.Now().Format(time.RFC3339), taskID)
	return err
}

func GetTaskByID(cfg *model.Config, taskID int) (*model.Task, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow(`
		SELECT id, description, deadline, status,
			COALESCE(completion_message, ''),
			COALESCE(review_message, ''),
			COALESCE(reviewed_by, ''),
			COALESCE(reviewed_at, ''),
			user, title, author, id_project, COALESCE(id_user, 0)
		FROM tasks
		WHERE id = ?
	`, taskID)

	var t model.Task
	if err := row.Scan(
		&t.ID,
		&t.Description,
		&t.Deadline,
		&t.Status,
		&t.CompletionMessage,
		&t.ReviewMessage,
		&t.ReviewedBy,
		&t.ReviewedAt,
		&t.User,
		&t.Title,
		&t.Author,
		&t.IdProject,
		&t.IdUser,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}
