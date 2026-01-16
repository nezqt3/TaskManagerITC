package services

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"backend/internal/config"
	"backend/internal/model"
	"backend/internal/notifications"
)

func CreateTask(task *model.Task) error {
	cfg := config.LoadConfig()
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return err
	}

	DB, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return err
	}
	defer DB.Close()

	if task.Status == "" {
		task.Status = "Новая"
	}

	if task.IdUser == 0 && task.User != "" {
		if user, err := GetUserByUsername(cfg, task.User); err == nil && user != nil {
			if user.TelegramID != "" {
				if telegramID, err := strconv.ParseInt(user.TelegramID, 10, 64); err == nil {
					task.IdUser = telegramID
				}
			}
		}
	}

	var message string
	result, err := DB.Exec(`
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
		return err
	}

	id, _ := result.LastInsertId()
	task.ID = int(id)
	task.IdUser = int64(task.IdUser)
	project, err := GetProjectByID(cfg, task.IdProject)
	projectTitle := ""
	if err == nil && project != nil {
		projectTitle = project.Title
	}

	deadlineStr := ""
	if task.Deadline != "" {
		deadlineTime, err := time.Parse("2006-01-02", task.Deadline)
		if err != nil {
			return fmt.Errorf("invalid deadline format: %v", err)
		}
		deadlineStr = deadlineTime.Format("02.01.2006")
	}

	message = fmt.Sprintf(
		"📌 Вам пришла новая задача:\n\n"+
			"Проект: %s\n"+
			"Задача: %s\n"+
			"Описание: %s\n\n"+
			"👤 Исполнитель: %s\n"+
			"✍️ Автор: %s\n"+
			"⏰ Дедлайн: %s\n"+
			"🆔 ID задачи: %d",
		projectTitle,
		task.Title,
		task.Description,
		task.User,
		task.Author,
		deadlineStr,
		task.ID,
	)

	if task.IdUser != 0 {
		notifications.SendTelegramNotification(cfg, task.IdUser, message)
	}

	return nil
}

func GetTasksByProjectID(projectID int) ([]model.Task, error) {
	cfg := config.LoadConfig()
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return nil, err
	}

	DB, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	rows, err := DB.Query(`
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
		log.Fatal(err)
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
			log.Fatal(err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tasks, err
}

func UpdateTask(cfg *model.Config, task *model.Task) error {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return err
	}

	DB, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return err
	}
	defer DB.Close()

	_, err = DB.Exec(`
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
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return err
	}

	DB, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return err
	}
	defer DB.Close()

	_, err = DB.Exec(`DELETE FROM tasks WHERE id = ?`, taskID)
	return err
}

func SubmitTaskCompletion(cfg *model.Config, taskID int, message string) error {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return err
	}

	DB, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return err
	}
	defer DB.Close()

	_, err = DB.Exec(`
		UPDATE tasks
		SET status = ?, completion_message = ?, review_message = '', reviewed_by = '', reviewed_at = ''
		WHERE id = ?
	`, "На проверке", message, taskID)
	return err
}

func ReviewTaskCompletion(cfg *model.Config, taskID int, approved bool, reviewer string, message string) error {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return err
	}

	DB, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return err
	}
	defer DB.Close()

	status := "Отклонена"
	if approved {
		status = "Выполнена"
	}

	_, err = DB.Exec(`
		UPDATE tasks
		SET status = ?, review_message = ?, reviewed_by = ?, reviewed_at = ?
		WHERE id = ?
	`, status, message, reviewer, time.Now().Format(time.RFC3339), taskID)
	return err
}

func GetTaskByID(cfg *model.Config, taskID int) (*model.Task, error) {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return nil, err
	}

	DB, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	row := DB.QueryRow(`
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
