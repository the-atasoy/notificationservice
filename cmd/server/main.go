package main

import (
	"log"

	"notificationservice/internal/config"
	//"notificationservice/internal/repository"
)

func main() {
	log.Println("Starting Notification Service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize MongoDB repository
	// mongoRepo, err := repository.NewMongoRepository(cfg.MongoDB.URI, cfg.MongoDB.Database)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to MongoDB: %v", err)
	// }

	log.Printf("Server starting on port: %s", cfg.Server.Port)
}
