package service

import (
	"errors"
	"strconv"
	"time"

	"backend/internal/model"
	"backend/pkg/jwt"
)

func TelegramAuth(req *model.AuthRequest, cfg *model.Config) (*model.AuthResponse, error) {
	if req.ID == 0 || req.Hash == "" {
		return nil, errors.New("invalid auth request")
	}

	ttl, err := time.ParseDuration(cfg.JWTTTL)
	if err != nil {
		ttl = 24 * time.Hour
	}

	token, err := jwt.GenerateToken(req.ID, cfg.JWTSecret, ttl)
	if err != nil {
		return nil, err
	}

	telegramID := strconv.FormatInt(req.ID, 10)
	csvUser, err := GetUserByTelegramID(cfg, telegramID)
	if err != nil {
		return nil, err
	}

	username := req.Username
	if username == "" {
		username = csvUser.Username
	}

	fullName := csvUser.FullName
	if fullName == "" {
		if req.FirstName != "" || req.LastName != "" {
			if req.LastName != "" {
				fullName = req.FirstName + " " + req.LastName
			} else {
				fullName = req.FirstName
			}
		}
	}

	profile := model.UserProfile{
		TelegramID:     telegramID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Username:       username,
		PhotoURL:       req.PhotoURL,
		FullName:       fullName,
		DateOfBirthday: csvUser.DateOfBirthday,
		NumberOfPhone:  csvUser.NumberOfPhone,
		Role:           csvUser.Role,
		MayToOpen:      csvUser.MayToOpen,
	}

	resp := &model.AuthResponse{
		JWT:     token,
		Profile: profile,
	}

	return resp, nil
}
