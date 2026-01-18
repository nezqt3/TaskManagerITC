package services

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"backend/internal/config"
	"backend/internal/model"
	"backend/internal/notifications"
	"backend/internal/repository"
)

func CreateTask(task *model.Task) error {
	cfg := config.LoadConfig()

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
	id, err := repository.CreateTask(cfg, task)
	if err != nil {
		return err
	}

	task.ID = id
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
	tasks, err := repository.GetTasksByProjectID(cfg, projectID)
	if err != nil {
		log.Fatal(err)
	}

	return tasks, err
}

func UpdateTask(cfg *model.Config, task *model.Task) error {
	return repository.UpdateTask(cfg, task)
}

func DeleteTask(cfg *model.Config, taskID int) error {
	return repository.DeleteTask(cfg, taskID)
}

func SubmitTaskCompletion(cfg *model.Config, taskID int, message string) error {
	task, err := GetTaskByID(cfg, taskID)
	if err != nil {
		return err
	}

	project, _ := GetProjectByID(cfg, task.IdProject)
	projectTitle := ""
	if project != nil {
		projectTitle = project.Title
	}

	author, err := repository.GetUserByFullName(cfg, task.Author)
	if err == nil && author != nil && author.TelegramID != "" {

		notifyMsg := fmt.Sprintf(
			"✅ Исполнитель отправил решение по задаче\n\n"+
				"Проект: %s\n"+
				"Задача: %s\n"+
				"Исполнитель: %s\n"+
				"Сообщение:\n%s\n\n"+
				"🆔 ID задачи: %d",
			projectTitle,
			task.Title,
			task.User,
			message,
			task.ID,
		)

		telegramID, _ := strconv.ParseInt(author.TelegramID, 10, 64)
		notifications.SendTelegramNotification(cfg, telegramID, notifyMsg)
	}

	// 💾 сохраняем решение
	return repository.SubmitTaskCompletion(cfg, taskID, message)
}

func ReviewTaskCompletion(cfg *model.Config, taskID int, approved bool, reviewer string, message string) error {
	task, err := GetTaskByID(cfg, taskID)
	if err != nil {
		return err
	}

	project, err := GetProjectByID(cfg, task.IdProject)
	projectTitle := ""
	if err == nil && project != nil {
		projectTitle = project.Title
	}

	deadlineStr := "—"
	if task.Deadline != "" {
		deadlineTime, err := time.Parse("2006-01-02", task.Deadline)
		if err != nil {
			return fmt.Errorf("invalid deadline format: %v", err)
		}
		deadlineStr = deadlineTime.Format("02.01.2006")
	}

	statusText := "❌ Задача отклонена"
	if approved {
		statusText = "✅ Задача принята"
	}

	notificationMessage := fmt.Sprintf(
		"%s\n\n"+
			"Проект: %s\n"+
			"Задача: %s\n"+
			"👤 Исполнитель: %s\n"+
			"✍️ Проверил: %s\n"+
			"⏰ Дедлайн: %s\n"+
			"💬 Комментарий:\n%s\n\n"+
			"🆔 ID задачи: %d",
		statusText,
		projectTitle,
		task.Title,
		task.User,
		reviewer,
		deadlineStr,
		message,
		task.ID,
	)

	if task.User != "" {
		user, err := GetUserByUsername(cfg, task.User)
		if err == nil && user != nil && user.TelegramID != "" {
			if telegramID, err := strconv.ParseInt(user.TelegramID, 10, 64); err == nil {
				notifications.SendTelegramNotification(cfg, telegramID, notificationMessage)
			}
		}
	}

	return repository.ReviewTaskCompletion(cfg, taskID, approved, reviewer, message)
}

func GetTaskByID(cfg *model.Config, taskID int) (*model.Task, error) {
	return repository.GetTaskByID(cfg, taskID)
}
