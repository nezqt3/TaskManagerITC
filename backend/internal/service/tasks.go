package service

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"path/filepath"
	"log"
	"fmt"
	"time"

	"backend/internal/model"
	"backend/internal/config"
	"backend/internal/notifications"
)

func CreateTask(task *model.Task) error {
	cfg := config.LoadConfig()
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return err
	}

	DB, err := sql.Open(cfg.DATABASE, absPath)
    var message string
	result, err := DB.Exec(`
		INSERT INTO tasks (
			description,
			deadline,
			status,
            id_user,
			user,
			title,
			author,
			id_project
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.Description,
		task.Deadline,
		task.Status,
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
    projectTitle := project.Title

    deadlineTime, err := time.Parse("2006-01-02", task.Deadline)
	if err != nil {
		return fmt.Errorf("invalid deadline format: %v", err)
	}

	deadlineStr := deadlineTime.Format("02.01.2006")

    if err != nil {
        fmt.Println("Ошибка %v \n", err)
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

    notifications.SendTelegramNotification(cfg, task.IdUser, message)

	return nil
}

func GetTasksByProjectID(projectID int) ([]model.Task, error) {
	cfg := config.LoadConfig()
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return nil, err
	}

	DB, err := sql.Open(cfg.DATABASE, absPath)
    rows, err := DB.Query("SELECT * FROM tasks WHERE id_project = ?", projectID)
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
            &t.User,
            &t.Title,
            &t.Author,
            &t.IdProject,
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