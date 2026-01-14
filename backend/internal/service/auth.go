package service

import (
	"errors"
	"time"

	"github.com/yourname/telegram-auth/internal/model"
	"github.com/yourname/telegram-auth/pkg/jwt"
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

	user := model.User{
		ID: req.ID,
		FirstName: req.FirstName,
		LastName: req.LastName,
		Username: req.Username,
		PhotoURL: req.PhotoURL,
	}

	resp := &model.AuthResponse{
		JWT: token,
		USER: user,
	}

	return resp, nil
}