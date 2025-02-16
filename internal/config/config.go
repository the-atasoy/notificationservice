package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    MongoDB struct {
        URI      string
        Database string
    }
    RabbitMQ struct {
        URI      string
        Queue    string
        Exchange string
    }
    Server struct {
        Port string
    }
}

func LoadConfig() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        return nil, err
    }

    config := &Config{}

    // MongoDB Config
    config.MongoDB.URI = getEnv("MONGODB_URI", "mongodb://admin:password@localhost:27017")
    config.MongoDB.Database = getEnv("MONGODB_DATABASE", "notifications")

    // RabbitMQ Config
    config.RabbitMQ.URI = getEnv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")
    config.RabbitMQ.Queue = getEnv("RABBITMQ_QUEUE", "notifications")
    config.RabbitMQ.Exchange = getEnv("RABBITMQ_EXCHANGE", "notifications")

    // Server Config
    config.Server.Port = getEnv("SERVER_PORT", "8080")

    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}