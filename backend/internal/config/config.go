package config 

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"backend/internal/model"
)

func LoadConfig() *model.Config{
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	cfg := &model.Config{
		AppPort: getEnv("APP_PORT", "8080"),
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		JWTSecret: getEnv("JWT_SECRET", "supersecret"),
		JWTTTL: getEnv("JWT_TTL", "24h"),
		DBDSN: getEnv("DBDSN", ""),
	}

	return cfg
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}