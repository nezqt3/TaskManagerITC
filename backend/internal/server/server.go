package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"backend/internal/handler"
	"backend/internal/model"
	"backend/internal/model/database"
	"backend/internal/service"
	"backend/internal/middleware"
)

type App struct {
	router http.Handler
	cfg    *model.Config
}

func WrapMiddleware(h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.Handler {
	var handler http.Handler = h
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}

func New(cfg *model.Config) *App {
	mux := http.NewServeMux()

	// проверка работы api
	mux.Handle("/health", WrapMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello"))
		}),
		middleware.JWTMiddleware(cfg.JWTSecret),
		middleware.RoleMiddleware("Владелец", "Руководитель", "Почётный"),
	))

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

	// end-point получения проектов
	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		var (
			projects []model.Project
			err      error
		)

		if id != "" {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
				return
			}
			projects, err = service.GetProjectsByID(cfg, idInt)
		} else {
			projects, err = service.GetProjects(cfg)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(projects)
	})

	// end-point получения одного проекта
	mux.HandleFunc("/projects/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idRaw := strings.TrimPrefix(r.URL.Path, "/projects/")
		if idRaw == "" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		project, err := service.GetProjectByID(cfg, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if project == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(project)
	})

	// end-point получение/создание тасок по задаче
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			id := r.URL.Query().Get("id_project")
			if id == "" {
				http.Error(w, "id_project is required", http.StatusBadRequest)
				return
			}

			idInt, err := strconv.Atoi(id)
			if err != nil {
				http.Error(w, "invalid id_project", http.StatusBadRequest)
				return
			}

			tasks := handler.GetTasksByProjectID(idInt)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)

		case http.MethodPost:
			var input database.Task

			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			task := database.Task{
				Description: input.Description,
				Deadline:    input.Deadline,
				Status:      input.Status,
				User:        input.User,
				Title:       input.Title,
				Author:      input.Author,
				IdProject:   input.IdProject,
				IdUser:      input.IdUser,
			}

			if err := handler.CreateTask(&task); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(task)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
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
