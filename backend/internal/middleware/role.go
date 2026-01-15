package middleware

import (
	"net/http"
	"strings"
)

func RoleMiddleware(allowedRoles ...string ) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value("role").(string)
			if !ok {
				http.Error(w, "unauthorized: role not found", http.StatusUnauthorized)
				return
			}

			for _, allowed := range allowedRoles {
				if strings.Contains(strings.ToLower(allowed), strings.ToLower(role)) {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "access denied", http.StatusForbidden)
		})
	}
}