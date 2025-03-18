package main

import (
	"context"
	"log"
	"time"

	"notificationservice/internal/config"
	"notificationservice/internal/repository"
	"notificationservice/internal/rabbitmq"
)

func main() {
    log.Println("Starting Notification Service...")

    // Create a background context
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cfg, err := config.LoadConfig()
    if (err != nil) {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize MongoDB repository
    mongoRepo, err := repository.NewMongoRepository(cfg.MongoDB.URI, cfg.MongoDB.Database)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    // Verify MongoDB connection
    if err := mongoRepo.Ping(ctx); err != nil {
        log.Fatalf("Failed to ping MongoDB: %v", err)
    }
    log.Println("Successfully connected to MongoDB")

    // Ensure MongoDB disconnects properly on application shutdown
    defer func() {
        if err := mongoRepo.Disconnect(context.Background()); err != nil {
            log.Printf("Error disconnecting from MongoDB: %v", err)
        }
    }()

    // Initialize RabbitMQ consumer
    consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQ.URI, cfg.RabbitMQ.Queue)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }

    // Start consuming messages
    err = consumer.Start(func(msg []byte) error {
        // TODO: Implement message handling
        log.Printf("Received message: %s", string(msg))
        return nil
    })
    if err != nil {
        log.Fatalf("Failed to start consuming messages: %v", err)
    }

    log.Printf("Server starting on port: %s", cfg.Server.Port)
}
