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
        RoutingKey string
        DeadLetterQueue struct {
            Queue      string
            Exchange   string
            RoutingKey string
        }
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
    config.RabbitMQ.RoutingKey = os.Getenv("RABBITMQ_ROUTING_KEY")
    config.RabbitMQ.DeadLetterQueue.Queue = os.Getenv("RABBITMQ_DLQ_QUEUE")
    config.RabbitMQ.DeadLetterQueue.Exchange = os.Getenv("RABBITMQ_DLQ_EXCHANGE")
    config.RabbitMQ.DeadLetterQueue.RoutingKey = os.Getenv("RABBITMQ_DLQ_ROUTING_KEY")

    config.Server.Port = os.Getenv("SERVER_PORT")

    return config, nil
}