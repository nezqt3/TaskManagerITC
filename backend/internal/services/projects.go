package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"backend/internal/model"
)

func GetProjects(cfg *model.Config) ([]model.Project, error) {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, description, status, users FROM projects ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var project model.Project
		var rawUsers string

		if err := rows.Scan(&project.ID, &project.Title, &project.Description, &project.Status, &rawUsers); err != nil {
			return nil, err
		}

		if rawUsers != "" {
			if err := json.Unmarshal([]byte(rawUsers), &project.Members); err != nil {
				return nil, fmt.Errorf("invalid project users: %w", err)
			}
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func GetProjectsByUsername(cfg *model.Config, username string) ([]model.Project, error) {
	normalized := normalizeUsername(username)
	if normalized == "" {
		return []model.Project{}, nil
	}

	projects, err := GetProjects(cfg)
	if err != nil {
		return nil, err
	}

	var filtered []model.Project
	for _, project := range projects {
		for _, member := range project.Members {
			if normalizeUsername(member.Username) == normalized {
				filtered = append(filtered, project)
				break
			}
		}
	}

	return filtered, nil
}

func UpdateProjectMembers(cfg *model.Config, projectID int, members []model.ProjectMember) error {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return err
	}

	db, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return err
	}
	defer db.Close()

	payload, err := json.Marshal(members)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE projects SET users = ? WHERE id = ?`, string(payload), projectID)
	return err
}

func UpdateProjectStatus(cfg *model.Config, projectID int, status string) error {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return err
	}

	db, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`UPDATE projects SET status = ? WHERE id = ?`, status, projectID)
	return err
}

func AddProjectMember(cfg *model.Config, projectID int, member model.ProjectMember) error {
	project, err := GetProjectByID(cfg, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}

	normalized := normalizeUsername(member.Username)
	for _, existing := range project.Members {
		if normalizeUsername(existing.Username) == normalized {
			return fmt.Errorf("member exists")
		}
	}
	project.Members = append(project.Members, member)
	return UpdateProjectMembers(cfg, projectID, project.Members)
}

func UpdateProjectMemberRole(cfg *model.Config, projectID int, username string, role string) error {
	project, err := GetProjectByID(cfg, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}

	normalized := normalizeUsername(username)
	updated := false
	for i, member := range project.Members {
		if normalizeUsername(member.Username) == normalized {
			project.Members[i].Role = role
			updated = true
			break
		}
	}
	if !updated {
		return fmt.Errorf("member not found")
	}
	return UpdateProjectMembers(cfg, projectID, project.Members)
}

func RemoveProjectMember(cfg *model.Config, projectID int, username string) error {
	project, err := GetProjectByID(cfg, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}

	normalized := normalizeUsername(username)
	filtered := make([]model.ProjectMember, 0, len(project.Members))
	for _, member := range project.Members {
		if normalizeUsername(member.Username) == normalized {
			continue
		}
		filtered = append(filtered, member)
	}
	if len(filtered) == len(project.Members) {
		return fmt.Errorf("member not found")
	}
	return UpdateProjectMembers(cfg, projectID, filtered)
}

func UpdateProjectStatusFromTasks(cfg *model.Config, projectID int) error {
	absPath, err := filepath.Abs(cfg.NAME_OF_DATABASE)
	if err != nil {
		return err
	}

	db, err := sql.Open(cfg.DATABASE, absPath)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT status FROM tasks WHERE id_project = ?`, projectID)
	if err != nil {
		return err
	}
	defer rows.Close()

	allCompleted := true
	hasTasks := false
	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err != nil {
			return err
		}
		hasTasks = true
		if strings.ToLower(status) != strings.ToLower("Выполнена") {
			allCompleted = false
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	project, err := GetProjectByID(cfg, projectID)
	if err != nil || project == nil {
		return err
	}

	if hasTasks && allCompleted {
		if project.Status != "Выполнен" {
			return UpdateProjectStatus(cfg, projectID, "Выполнен")
		}
		return nil
	}

	if project.Status == "Выполнен" {
		return UpdateProjectStatus(cfg, projectID, "В работе")
	}

	return nil
}

func GetProjectsByID(cfg *model.Config, id int) ([]model.Project, error) {
	project, err := GetProjectByID(cfg, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return []model.Project{}, nil
	}
	return []model.Project{*project}, nil
}

func GetProjectByID(cfg *model.Config, id int) (*model.Project, error) {
	projects, err := GetProjects(cfg)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.ID == id {
			return &project, nil
		}
	}

	return nil, nil
}

func normalizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.TrimPrefix(username, "@")
	return strings.ToLower(username)
}
