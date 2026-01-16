package services

import (
	"backend/internal/logger"
	"backend/internal/model"
	"backend/internal/repository"
)

func GetEvents(cfg *model.Config) ([]model.Event, error) {
	events, err := repository.GetEvents(cfg)
	if err != nil {
		logger.Error.Printf("GetEvents: failed to load events: %v\n", err)
		return nil, err
	}

	return events, nil
}
