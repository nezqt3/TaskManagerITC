package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "unauthorized: missing authorization", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims := &model.Claims{}

			_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
			ctx = context.WithValue(ctx, "role", claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}