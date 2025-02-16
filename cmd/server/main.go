package main

import (
	"log"

	"github.com/the-atasoy/notificationservice/internal/config"
)

func main() {
	log.Println("Starting Notification Service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Server starting on port: %s", cfg.Server.Port)
}
