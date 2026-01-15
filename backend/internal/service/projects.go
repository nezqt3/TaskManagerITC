package service

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
