package server

import (
	"net/http"
	
	"github.com/yourname/telegram-auth/internal/config"
	"github.com/yourname/telegram-auth/internal/model"
)

type App struct {
	router http.Handler
	cfg *model.Config
}

func New(cfg *model.Config) *App {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello"))
	})

	return &App {
		router: mux,
		cfg: cfg,
	}

}

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, a.router)
}