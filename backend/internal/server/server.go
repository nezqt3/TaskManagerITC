package server

import (
	"encoding/json"
	"net/http"

	"backend/internal/handler"
	"backend/internal/model"
	"backend/internal/service"
)

type App struct {
	router http.Handler
	cfg    *model.Config
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

	// end-point получения пользователей
	mux.HandleFunc("/get_users", func(w http.ResponseWriter, r *http.Request) {
		users, err := service.GetUsers(cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	return &App{
		router: withCORS(mux),
		cfg:    cfg,
	}

}

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, a.router)
}

func withCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
