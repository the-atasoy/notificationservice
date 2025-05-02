package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"notificationservice/internal/config"
	"notificationservice/internal/handlers"
	"notificationservice/internal/rabbitmq"
	"notificationservice/internal/repository"
)

func main() {
    log.Println("Starting Notification Service...")

    // Load config
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    mongoRepo, err := repository.NewMongoRepository(cfg.MongoDB.URI, cfg.MongoDB.Database)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    notificationHandler := handlers.NewNotificationHandler(mongoRepo)

    consumer := rabbitmq.NewConsumer(
        cfg.RabbitMQ.URI,
        cfg.RabbitMQ.Queue,
        cfg.RabbitMQ.Exchange,
        cfg.RabbitMQ.RoutingKey,
        &rabbitmq.ConsumerOptions{
            DeadLetterConfig: &rabbitmq.DeadLetterConfig{
                ExchangeName: cfg.RabbitMQ.DeadLetterQueue.Exchange,
                QueueName:    cfg.RabbitMQ.DeadLetterQueue.Queue,
                RoutingKey:   cfg.RabbitMQ.DeadLetterQueue.RoutingKey,
            },
        },
    )

    if err := consumer.Connect(); err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer consumer.Close()

    if err := consumer.Start(notificationHandler); err != nil {
        log.Fatalf("Failed to start consuming messages: %v", err)
    }

    log.Printf("Server started successfully")

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")
}
