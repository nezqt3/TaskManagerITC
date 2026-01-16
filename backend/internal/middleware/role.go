package middleware

import (
	"net/http"
	"strings"

	"backend/internal/logger"
)

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value("role").(string)
			if !ok {
				logger.Error.Println("RoleMiddleware: role not found in context")
				http.Error(w, "unauthorized: role not found", http.StatusUnauthorized)
				return
			}

			for _, allowed := range allowedRoles {
				if strings.EqualFold(allowed, role) {
					next.ServeHTTP(w, r)
					return
				}
			}

			logger.Error.Printf("RoleMiddleware: access denied for role '%s'\n", role)
			http.Error(w, "access denied", http.StatusForbidden)
		})
	}
}
