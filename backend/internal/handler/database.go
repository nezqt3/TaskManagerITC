package handler

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"path/filepath"
	"log"
    "fmt"
    "time"

    "backend/internal/notifications"
	"backend/internal/model"
    "backend/internal/config"
    "backend/internal/service"
	"backend/internal/model/database"
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

	taskTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id    INTEGER PRIMARY KEY, 
		description TEXT, 
		deadline DATE, 
		status TEXT, 
        user_id INTEGER,
		user TEXT, 
		title TEXT,
		author TEXT,
		id_project INTEGER
	);
	`
	_, err = DB.Exec(taskTable)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTask(task *database.Task) error {
    var message string
    cfg := config.LoadConfig()
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
    project, err := service.GetProjectByID(cfg, task.IdProject)
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

func GetTasksByProjectID(projectID int) []database.Task {
    rows, err := DB.Query("SELECT * FROM tasks WHERE id_project = ?", projectID)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var tasks []database.Task
    for rows.Next() {
        var t database.Task
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

    return tasks
}