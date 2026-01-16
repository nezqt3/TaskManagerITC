package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/logger"
	"backend/internal/server"
)

func main() {
	cfg := config.LoadConfig()

	if err := logger.Init("logs/app.log"); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	logger.Info.Println("application starting")

	handler.InitDatabase(cfg)
	database.RunMigrations()
	app := server.New(cfg)

	logger.Info.Printf("server started on port: %s", cfg.AppPort)

	if err := app.Run(":" + cfg.AppPort); err != nil {
		logger.Fatal.Println("server stopped with error:", err)
	}
}
