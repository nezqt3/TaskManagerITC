package services

import (
	"errors"
	"strconv"
	"time"

	"backend/internal/logger"
	"backend/internal/model"
	"backend/pkg/jwt"
)

func TelegramAuth(req *model.AuthRequest, cfg *model.Config) (*model.AuthResponse, error) {
	if req.ID == 0 || req.Hash == "" {
		logger.Error.Println("TelegramAuth: invalid auth request - missing ID or Hash")
		return nil, errors.New("invalid auth request")
	}

	telegramID := strconv.FormatInt(req.ID, 10)
	dbUser, err := GetUserByTelegramID(cfg, telegramID)
	if err != nil {
		logger.Error.Printf("TelegramAuth: failed to get user by TelegramID '%s': %v\n", telegramID, err)
		return nil, err
	}

	username := req.Username
	if username == "" {
		username = dbUser.Username
	}

	fullName := dbUser.FullName
	if fullName == "" {
		if req.FirstName != "" || req.LastName != "" {
			if req.LastName != "" {
				fullName = req.FirstName + " " + req.LastName
			} else {
				fullName = req.FirstName
			}
		}
	}

	firstName := req.FirstName
	if firstName == "" {
		firstName = dbUser.FirstName
	}

	lastName := req.LastName
	if lastName == "" {
		lastName = dbUser.LastName
	}

	profile := model.UserProfile{
		TelegramID:     telegramID,
		FirstName:      firstName,
		LastName:       lastName,
		Username:       username,
		PhotoURL:       req.PhotoURL,
		FullName:       fullName,
		DateOfBirthday: dbUser.DateOfBirthday,
		NumberOfPhone:  dbUser.NumberOfPhone,
		Role:           dbUser.Role,
		MayToOpen:      dbUser.MayToOpen,
	}

	ttl, err := time.ParseDuration(cfg.JWTTTL)
	if err != nil {
		logger.Error.Printf("TelegramAuth: invalid JWTTTL '%s', defaulting to 24h: %v\n", cfg.JWTTTL, err)
		ttl = 24 * time.Hour
	}

	token, err := jwt.GenerateToken(req.ID, dbUser.Role, cfg.JWTSecret, ttl)
	if err != nil {
		logger.Error.Printf("TelegramAuth: failed to generate JWT for TelegramID '%s': %v\n", telegramID, err)
		return nil, err
	}

	logger.Info.Printf("TelegramAuth: successful login for TelegramID '%s', username '%s'\n", telegramID, username)

	resp := &model.AuthResponse{
		JWT:     token,
		Profile: profile,
	}

	return resp, nil
}
