package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/model"
	"backend/internal/permissions"
	"backend/internal/services"
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
	mux.Handle("/get_users", WrapMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			users, err := services.GetUsers(cfg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)
		}),
		middleware.JWTMiddleware(cfg.JWTSecret),
	))

	// end-point обновления пользователя
	mux.Handle("/users/", WrapMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			role, _ := r.Context().Value("role").(string)
			userID, _ := r.Context().Value("user_id").(int64)
			if !permissions.IsAdmin(role) {
				user, err := services.GetUserByTelegramID(cfg, strconv.FormatInt(userID, 10))
				if err != nil || user == nil || !permissions.IsAdmin(user.Role) {
					http.Error(w, "access denied", http.StatusForbidden)
					return
				}
			}

			telegramID := strings.TrimPrefix(r.URL.Path, "/users/")
			if telegramID == "" {
				http.Error(w, "telegram id is required", http.StatusBadRequest)
				return
			}

			var payload model.UserProfile
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			payload.Username = strings.TrimPrefix(strings.TrimSpace(payload.Username), "@")
			if err := services.UpdateUser(cfg, telegramID, &payload); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}),
		middleware.JWTMiddleware(cfg.JWTSecret),
	))

	// end-point получения проектов
	mux.Handle("/projects", WrapMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			id := r.URL.Query().Get("id")
			username := r.URL.Query().Get("username")
			var (
				projects []model.Project
				err      error
			)

			if username != "" {
				projects, err = services.GetProjectsByUsername(cfg, username)
			} else if id != "" {
				idInt, err := strconv.Atoi(id)
				if err != nil {
					http.Error(w, "invalid id", http.StatusBadRequest)
					return
				}
				projects, err = services.GetProjectsByID(cfg, idInt)
			} else {
				projects, err = services.GetProjects(cfg)
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(projects)
		}),
		middleware.JWTMiddleware(cfg.JWTSecret),
	))

	// end-point получения/управления одним проектом
	mux.Handle("/projects/", WrapMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/projects/")
			parts := strings.Split(strings.Trim(path, "/"), "/")
			if len(parts) == 0 || parts[0] == "" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}

			id, err := strconv.Atoi(parts[0])
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
				return
			}

			if len(parts) == 1 {
				if r.Method != http.MethodGet {
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					return
				}

				project, err := services.GetProjectByID(cfg, id)
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
				return
			}

			sub := parts[1]
			switch sub {
			case "status":
				if r.Method != http.MethodPut {
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					return
				}

				role, _ := r.Context().Value("role").(string)
				userID, _ := r.Context().Value("user_id").(int64)
				if !canManageProjectTasks(cfg, id, userID, role) {
					http.Error(w, "access denied", http.StatusForbidden)
					return
				}

				var payload struct {
					Status string `json:"status"`
				}
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					http.Error(w, "invalid json", http.StatusBadRequest)
					return
				}
				status := strings.TrimSpace(payload.Status)
				if status == "" {
					http.Error(w, "status is required", http.StatusBadRequest)
					return
				}

				if err := services.UpdateProjectStatus(cfg, id, status); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			case "members":
				role, _ := r.Context().Value("role").(string)
				userID, _ := r.Context().Value("user_id").(int64)
				if !permissions.IsAdmin(resolveEffectiveRole(cfg, userID, role)) {
					http.Error(w, "access denied", http.StatusForbidden)
					return
				}

				switch r.Method {
				case http.MethodPost:
					var payload model.ProjectMember
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
						http.Error(w, "invalid json", http.StatusBadRequest)
						return
					}
					payload.Username = strings.TrimSpace(strings.TrimPrefix(payload.Username, "@"))
					if payload.Username == "" {
						http.Error(w, "username is required", http.StatusBadRequest)
						return
					}
					if invalidRoleCombo(payload.Role) {
						http.Error(w, "invalid role combination", http.StatusBadRequest)
						return
					}

					if payload.FullName == "" {
						if user, err := services.GetUserByUsername(cfg, payload.Username); err == nil && user != nil {
							payload.FullName = user.FullName
							payload.TelegramID = user.TelegramID
						}
					}

					payload.Role = normalizeMemberRole(payload.Role)
					if err := services.AddProjectMember(cfg, id, payload); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					w.WriteHeader(http.StatusCreated)
					return
				case http.MethodPut:
					if len(parts) < 3 {
						http.Error(w, "username is required", http.StatusBadRequest)
						return
					}
					username := parts[2]
					var payload struct {
						Role string `json:"role"`
					}
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
						http.Error(w, "invalid json", http.StatusBadRequest)
						return
					}
					if invalidRoleCombo(payload.Role) {
						http.Error(w, "invalid role combination", http.StatusBadRequest)
						return
					}
					roleValue := normalizeMemberRole(payload.Role)
					if err := services.UpdateProjectMemberRole(cfg, id, username, roleValue); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					w.WriteHeader(http.StatusNoContent)
					return
				case http.MethodDelete:
					if len(parts) < 3 {
						http.Error(w, "username is required", http.StatusBadRequest)
						return
					}
					username := parts[2]
					if err := services.RemoveProjectMember(cfg, id, username); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					w.WriteHeader(http.StatusNoContent)
					return
				default:
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					return
				}
			default:
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
		}),
		middleware.JWTMiddleware(cfg.JWTSecret),
	))

	// end-point получение/создание тасок по задаче
	mux.Handle("/tasks", WrapMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

				tasks, _ := services.GetTasksByProjectID(idInt)

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tasks)

			case http.MethodPost:
				role, _ := r.Context().Value("role").(string)
				userID, _ := r.Context().Value("user_id").(int64)

				var input model.Task

				if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
					http.Error(w, "invalid json", http.StatusBadRequest)
					return
				}

				task := model.Task{
					Description: input.Description,
					Deadline:    input.Deadline,
					Status:      input.Status,
					User:        input.User,
					Title:       input.Title,
					Author:      input.Author,
					IdProject:   input.IdProject,
					IdUser:      input.IdUser,
				}

				if !canManageProjectTasks(cfg, task.IdProject, userID, role) {
					http.Error(w, "access denied", http.StatusForbidden)
					return
				}

				if err := services.CreateTask(&task); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(task)

			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}),
		middleware.JWTMiddleware(cfg.JWTSecret),
	))

	// end-point управления задачами
	mux.Handle("/tasks/", WrapMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/tasks/")
			parts := strings.Split(strings.Trim(path, "/"), "/")
			if len(parts) == 0 || parts[0] == "" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}

			id, err := strconv.Atoi(parts[0])
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
				return
			}

			if len(parts) == 1 {
				switch r.Method {
				case http.MethodPut:
					role, _ := r.Context().Value("role").(string)
					userID, _ := r.Context().Value("user_id").(int64)

					var payload model.Task
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
						http.Error(w, "invalid json", http.StatusBadRequest)
						return
					}

					payload.ID = id
					if payload.IdProject == 0 {
						if existing, err := services.GetTaskByID(cfg, id); err == nil && existing != nil {
							payload.IdProject = existing.IdProject
						}
					}
					if !canManageProjectTasks(cfg, payload.IdProject, userID, role) {
						http.Error(w, "access denied", http.StatusForbidden)
						return
					}
					if payload.IdUser == 0 && payload.User != "" {
						if user, err := services.GetUserByUsername(cfg, payload.User); err == nil && user != nil {
							if user.TelegramID != "" {
								if telegramID, err := strconv.ParseInt(user.TelegramID, 10, 64); err == nil {
									payload.IdUser = telegramID
								}
							}
						}
					}

					if err := services.UpdateTask(cfg, &payload); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusNoContent)
					return
				case http.MethodDelete:
					role, _ := r.Context().Value("role").(string)
					userID, _ := r.Context().Value("user_id").(int64)

					var projectID int
					if existing, err := services.GetTaskByID(cfg, id); err == nil && existing != nil {
						projectID = existing.IdProject
					}
					if !canManageProjectTasks(cfg, projectID, userID, role) {
						http.Error(w, "access denied", http.StatusForbidden)
						return
					}

					if err := services.DeleteTask(cfg, id); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusNoContent)
					return
				default:
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					return
				}
			}

			sub := parts[1]
			switch sub {
			case "complete":
				if r.Method != http.MethodPost {
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					return
				}

				task, err := services.GetTaskByID(cfg, id)
				if err != nil || task == nil {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}

				userID, _ := r.Context().Value("user_id").(int64)
				role, _ := r.Context().Value("role").(string)
				if !canSubmitCompletion(cfg, task, userID, role) {
					http.Error(w, "access denied", http.StatusForbidden)
					return
				}

				var payload struct {
					Message string `json:"message"`
				}
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					http.Error(w, "invalid json", http.StatusBadRequest)
					return
				}

				if err := services.SubmitTaskCompletion(cfg, id, strings.TrimSpace(payload.Message)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusNoContent)
				return
			case "review":
				if r.Method != http.MethodPost {
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					return
				}

				role, _ := r.Context().Value("role").(string)
				userID, _ := r.Context().Value("user_id").(int64)

				var payload struct {
					Approved bool   `json:"approved"`
					Message  string `json:"message"`
				}
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					http.Error(w, "invalid json", http.StatusBadRequest)
					return
				}

				reviewer := resolveReviewerName(cfg, userID)

				task, err := services.GetTaskByID(cfg, id)
				if err != nil || task == nil {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}

				if !canReviewProjectTasks(cfg, task.IdProject, userID, role) {
					http.Error(w, "access denied", http.StatusForbidden)
					return
				}

				if err := services.ReviewTaskCompletion(cfg, id, payload.Approved, reviewer, strings.TrimSpace(payload.Message)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusNoContent)
				return
			default:
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
		}),
		middleware.JWTMiddleware(cfg.JWTSecret),
	))

	// dashboard data
	mux.Handle("/dashboard", WrapMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			username := r.URL.Query().Get("username")
			if username == "" {
				http.Error(w, "username is required", http.StatusBadRequest)
				return
			}

			projects, err := services.GetProjectsByUsername(cfg, username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			projectMap := make(map[int]string)
			for _, project := range projects {
				projectMap[project.ID] = project.Title
			}

			normalizedUsername := normalizeUsername(username)
			tasks := make([]model.Task, 0)
			for _, project := range projects {
				projectTasks, _ := services.GetTasksByProjectID(project.ID)
				for _, task := range projectTasks {
					if strings.ToLower(task.Status) == strings.ToLower("Выполнена") {
						continue
					}
					if normalizedUsername != "" &&
						normalizeUsername(task.User) != normalizedUsername {
						continue
					}
					task.ProjectTitle = projectMap[task.IdProject]
					tasks = append(tasks, task)
				}
			}

			events, err := services.GetEvents(cfg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			payload := model.Dashboard{
				Projects: projects,
				Tasks:    tasks,
				Events:   events,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(payload)
		}),
		middleware.JWTMiddleware(cfg.JWTSecret),
	))

	// events list
	mux.Handle("/events", WrapMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			events, err := services.GetEvents(cfg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(events)
		}),
		middleware.JWTMiddleware(cfg.JWTSecret),
	))

	mux.Handle("/search_users", WrapMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			term := strings.TrimSpace(r.URL.Query().Get("term"))
			if term == "" {
				http.Error(w, "term is required", http.StatusBadRequest)
				return
			}

			users, err := services.SearchUsersByFullName(cfg, term)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)
		}),
		middleware.JWTMiddleware(cfg.JWTSecret),
	))

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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func resolveReviewerName(cfg *model.Config, userID int64) string {
	if userID == 0 {
		return ""
	}

	user, err := services.GetUserByTelegramID(cfg, strconv.FormatInt(userID, 10))
	if err != nil || user == nil {
		return ""
	}

	if user.FullName != "" {
		return user.FullName
	}
	if user.FirstName != "" || user.LastName != "" {
		return strings.TrimSpace(user.FirstName + " " + user.LastName)
	}
	return user.Username
}

func resolveEffectiveRole(cfg *model.Config, userID int64, fallback string) string {
	if userID == 0 {
		return fallback
	}
	user, err := services.GetUserByTelegramID(cfg, strconv.FormatInt(userID, 10))
	if err != nil || user == nil || user.Role == "" {
		return fallback
	}
	return user.Role
}

func canManageProjectTasks(cfg *model.Config, projectID int, userID int64, role string) bool {
	effectiveRole := resolveEffectiveRole(cfg, userID, role)

	if permissions.IsAdmin(effectiveRole) {
		return true
	}
	if userID == 0 || projectID == 0 {
		return false
	}

	memberRole, isMember := getProjectMemberRole(cfg, projectID, userID)
	if !isMember {
		return false
	}

	if permissions.IsModerator(effectiveRole) {
		return true
	}

	return isLeaderRole(memberRole)
}

func canReviewProjectTasks(cfg *model.Config, projectID int, userID int64, role string) bool {
	return canManageProjectTasks(cfg, projectID, userID, role)
}

func getProjectMemberRole(cfg *model.Config, projectID int, userID int64) (string, bool) {
	project, err := services.GetProjectByID(cfg, projectID)
	if err != nil || project == nil {
		return "", false
	}

	user, err := services.GetUserByTelegramID(cfg, strconv.FormatInt(userID, 10))
	if err != nil || user == nil {
		return "", false
	}

	normalizedUsername := normalizeUsername(user.Username)
	for _, member := range project.Members {
		if member.TelegramID != "" && member.TelegramID == user.TelegramID {
			return member.Role, true
		}
		if normalizedUsername != "" && normalizeUsername(member.Username) == normalizedUsername {
			return member.Role, true
		}
	}

	return "", false
}

func normalizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.TrimPrefix(username, "@")
	return strings.ToLower(username)
}

func isLeaderRole(role string) bool {
	roles := permissions.ParseRoles(role)
	return roles["руководитель"]
}

func canSubmitCompletion(cfg *model.Config, task *model.Task, userID int64, role string) bool {
	return userID != 0
}

func normalizeMemberRole(role string) string {
	role = strings.TrimSpace(role)
	if role == "" {
		return role
	}

	if permissions.MustIncludeModerator(role) && !permissions.IsModerator(role) {
		return role + ", Модератор"
	}
	return role
}

func invalidRoleCombo(role string) bool {
	roles := permissions.ParseRoles(role)
	return roles["разработчик"] && (roles["админ"] || roles["admin"])
}
