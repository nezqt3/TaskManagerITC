package server

import (
	"net/http"
	"encoding/json"
	
	"backend/internal/model"
	"backend/internal/handler"
	"backend/internal/service"
)

type App struct {
	router http.Handler
	cfg *model.Config
}

func New(cfg *model.Config) *App {
	mux := http.NewServeMux()

	// проверка работы api
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello"))
	})

	// end-point авторизации
	mux.HandleFunc("/auth/telegram", handler.TelegramAuthHandler(cfg))

	mux.HandleFunc("/get_users", func(w http.ResponseWriter, r *http.Request) {
		users, err := service.GetUsers(cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	return &App {
		router: mux,
		cfg: cfg,
	}

}

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, a.router)
}