package main

import (
	"log"

	"backend/internal/server"
	"backend/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	app := server.New(cfg)

	log.Println("Server have already started on port", cfg.AppPort)
	if err := app.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}