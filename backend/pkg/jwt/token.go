package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"backend/internal/model"
)

func GenerateToken(userID int64, role string, secret string, ttl time.Duration) (string, error) {
	claims := model.Claims{
		UserId: userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte (secret))
}