package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"backend/internal/model"
	"backend/internal/services"
	"backend/internal/telegram"
)

func TelegramAuthHandler(cfg *model.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		dataMap := map[string]string{
			"id":         fmt.Sprintf("%d", req.ID),
			"first_name": req.FirstName,
			"auth_date":  fmt.Sprintf("%d", req.AuthDate),
			"hash":       req.Hash,
		}
		if req.LastName != "" {
			dataMap["last_name"] = req.LastName
		}
		if req.Username != "" {
			dataMap["username"] = req.Username
		}
		if req.PhotoURL != "" {
			dataMap["photo_url"] = req.PhotoURL
		}

		if err := telegram.CheckTelegramAuth(dataMap, cfg.TelegramBotToken); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		resp, err := services.TelegramAuth(&req, cfg)
		if err != nil {
			if errors.Is(err, services.ErrUserNotFound) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
