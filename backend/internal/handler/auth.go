package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"backend/internal/model"
	"backend/internal/services"
	"backend/internal/telegram"
	"backend/internal/logger"
)

func TelegramAuthHandler(cfg *model.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("TelegramAuthHandler called")
		
		var req model.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			logger.Error.Println(w, "invalid request", http.StatusBadRequest)
			return
		}
		logger.Info.Printf("Request decoded: %+v\n", req)

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

		logger.Info.Printf("Data map for Telegram auth: %+v\n", dataMap)

		if err := telegram.CheckTelegramAuth(dataMap, cfg.TelegramBotToken); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			logger.Error.Println(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		logger.Info.Println("Telegram auth successful")

		resp, err := services.TelegramAuth(&req, cfg)
		if err != nil {
			if errors.Is(err, services.ErrUserNotFound) {
				http.Error(w, "user not found", http.StatusNotFound)
				logger.Error.Println(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			logger.Error.Println(w, "internal error", http.StatusInternalServerError)
			return
		}
		logger.Info.Printf("Response sent successfully: %+v\n", resp)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
