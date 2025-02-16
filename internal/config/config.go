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

    config.MongoDB.URI = os.Getenv("MONGODB_URI")
    config.MongoDB.Database = os.Getenv("MONGODB_DATABASE")

    config.RabbitMQ.URI = os.Getenv("RABBITMQ_URI")
    config.RabbitMQ.Queue = os.Getenv("RABBITMQ_QUEUE")
    config.RabbitMQ.Exchange = os.Getenv("RABBITMQ_EXCHANGE")

    config.Server.Port = os.Getenv("SERVER_PORT")

    return config, nil
}