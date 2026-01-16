package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/model"
	"backend/internal/repository"
)

func GetProjects(cfg *model.Config) ([]model.Project, error) {
	rows, err := repository.GetProjects(cfg)
	if err != nil {
		return nil, err
	}

	var projects []model.Project
	for _, row := range rows {
		project := model.Project{
			ID:          row.ID,
			Title:       row.Title,
			Description: row.Description,
			Status:      row.Status,
		}

		rawUsers := row.RawUsers
		if rawUsers != "" {
			if err := json.Unmarshal([]byte(rawUsers), &project.Members); err != nil {
				return nil, fmt.Errorf("invalid project users: %w", err)
			}
		}

		projects = append(projects, project)
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
	payload, err := json.Marshal(members)
	if err != nil {
		return err
	}
	return repository.UpdateProjectMembers(cfg, projectID, string(payload))
}

func UpdateProjectStatus(cfg *model.Config, projectID int, status string) error {
	return repository.UpdateProjectStatus(cfg, projectID, status)
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
	statuses, err := repository.GetTaskStatusesByProjectID(cfg, projectID)
	if err != nil {
		return err
	}

	allCompleted := true
	hasTasks := len(statuses) > 0
	for _, status := range statuses {
		if strings.ToLower(status) != strings.ToLower("Выполнена") {
			allCompleted = false
		}
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
