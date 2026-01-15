package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/server"
)

func main() {
	cfg := config.LoadConfig()
	handler.InitDatabase(cfg)
	app := server.New(cfg)

	log.Println("Server have already started on port", cfg.AppPort)
	if err := app.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
