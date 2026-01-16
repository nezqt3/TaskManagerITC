package services

import (
	"database/sql"
	"errors"
	"strings"

	"backend/internal/model"
	"backend/internal/repository"
)

var ErrUserNotFound = errors.New("user not found")

func GetUsers(cfg *model.Config) ([]model.UserProfile, error) {
	return repository.GetUsers(cfg)
}

func GetUserByTelegramID(cfg *model.Config, telegramID string) (*model.UserProfile, error) {
	normalizedID := strings.TrimSpace(telegramID)
	u, err := repository.GetUserByTelegramID(cfg, normalizedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func GetUserByUsername(cfg *model.Config, username string) (*model.UserProfile, error) {
	normalized := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(username)), "@")
	if normalized == "" {
		return nil, ErrUserNotFound
	}
	normalizedWithAt := "@" + normalized

	u, err := repository.GetUserByUsername(cfg, normalized, normalizedWithAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func SearchUsersByFullName(cfg *model.Config, fullName string) ([]model.UserProfile, error) {
	return repository.SearchUsersByFullName(cfg, fullName)
}

func UpdateUser(cfg *model.Config, telegramID string, updates *model.UserProfile) error {
	return repository.UpdateUser(cfg, telegramID, updates)
}
