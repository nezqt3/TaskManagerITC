package server

import (
	"net/http"
	
	"backend/internal/model"
	"backend/internal/handler"
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

	return &App {
		router: mux,
		cfg: cfg,
	}

}

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, a.router)
}