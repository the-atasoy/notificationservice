package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"notificationservice/internal/errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler interface {
    ProcessMessage([]byte) error
}

type Consumer struct {
    uri              string
    queueName        string
    exchangeName     string
    routingKey       string
    deadLetterConfig *DeadLetterConfig
    connection       *amqp.Connection
    channel          *amqp.Channel
}

type DeadLetterConfig struct {
    QueueName    string
    ExchangeName string
    RoutingKey   string
}

type ConsumerOptions struct {
    DeadLetterConfig *DeadLetterConfig
}

func NewConsumer(uri, queueName, exchangeName string, routingKey string, options *ConsumerOptions) *Consumer {
    return &Consumer{
        uri:          uri,
        queueName:    queueName,
        exchangeName: exchangeName,
        routingKey: routingKey,
        deadLetterConfig: options.DeadLetterConfig,
    }
}

func (consumer *Consumer) Connect() error {
    connection, err := amqp.Dial(consumer.uri)
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }
    consumer.connection = connection

    channel, err := connection.Channel()
    if err != nil {
        return fmt.Errorf("failed to open channel: %w", err)
    }
    consumer.channel = channel

    err = consumer.channel.ExchangeDeclare(
        consumer.exchangeName, // name
        "direct",       // type
        true,           // durable
        false,          // auto-deleted
        false,          // internal
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare exchange: %w", err)
    }

    if consumer.deadLetterConfig != nil {
        err = consumer.setupDeadLetterExchange()
        if err != nil {
            return err
        }
    }

    args := amqp.Table{}
    if consumer.deadLetterConfig != nil {
        args["x-dead-letter-exchange"] = consumer.deadLetterConfig.ExchangeName
        args["x-dead-letter-routing-key"] = consumer.deadLetterConfig.RoutingKey
    }

    queue, err := consumer.channel.QueueDeclare(
        consumer.queueName, // name
        true,        // durable
        false,       // auto-delete
        false,       // exclusive
        false,       // no-wait
        args,        // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare queue: %w", err)
    }

    err = consumer.channel.QueueBind(
        queue.Name,
        consumer.routingKey,
        consumer.exchangeName,
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to bind queue: %w", err)
    }

    return nil
}

func (consumer *Consumer) setupDeadLetterExchange() error {
    err := consumer.channel.ExchangeDeclare(
        consumer.deadLetterConfig.ExchangeName, // name
        "direct",                        // type
        true,                            // durable
        false,                           // auto-deleted
        false,                           // internal
        false,                           // no-wait
        nil,                             // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare dead letter exchange: %w", err)
    }

    _, err = consumer.channel.QueueDeclare(
        consumer.deadLetterConfig.QueueName, // name
        true,                         // durable
        false,                        // auto-delete
        false,                        // exclusive
        false,                        // no-wait
        nil,                          // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare dead letter queue: %w", err)
    }

    err = consumer.channel.QueueBind(
        consumer.deadLetterConfig.QueueName,
        consumer.deadLetterConfig.RoutingKey,
        consumer.deadLetterConfig.ExchangeName,
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to bind dead letter queue: %w", err)
    }

    return nil
}

func (c *Consumer) Start(handler MessageHandler) error {
    if c.channel == nil {
        return fmt.Errorf("channel not initialized, call Connect() first")
    }

    msgs, err := c.channel.Consume(
        c.queueName, // queue
        "",          // consumer
        false,       // auto-ack
        false,       // exclusive
        true,        // no-local
        false,       // no-wait
        nil,         // args
    )
    if err != nil {
        return fmt.Errorf("failed to register consumer: %w", err)
    }

    for msg := range msgs {
        go func(msg amqp.Delivery) {
            log.Printf("Received message: %s", string(msg.Body))
            err := handler.ProcessMessage(msg.Body)
            if err != nil {
                log.Printf("Error processing message: %v", err)
                
                if errors.IsValidationError(err) {
                    log.Printf("Validation error detected, sending to dead letter queue")
                    
                    errorDesc := errors.GetErrorDescription(err)
                    c.moveToDeadLetter(msg.Body, string(errors.ValidationError), errorDesc)
                    
                    msg.Ack(false)
                } else if errors.IsRetriableError(err) {
                    log.Printf("Retriable error detected, requeueing message")
                    msg.Nack(false, true)
                } else {
                    log.Printf("Processing error detected, sending to dead letter queue")
                    
                    errorDesc := errors.GetErrorDescription(err)
                    c.moveToDeadLetter(msg.Body, string(errors.ProcessingError), errorDesc)
                    
                    msg.Ack(false)
                }
            } else {
                log.Printf("Message processed successfully")
                msg.Ack(false)
            }
        }(msg)
    }

    log.Println("RabbitMQ consumer started successfully")
    return nil
}

func (c *Consumer) moveToDeadLetter(body []byte, errorType, errorMsg string) {
    if c.deadLetterConfig == nil {
        log.Println("Dead letter exchange not configured, discarding failed message")
        return
    }

    err := c.channel.Publish(
        c.deadLetterConfig.ExchangeName,  // exchange
        c.deadLetterConfig.RoutingKey,    // routing key
        false,                            // mandatory
        false,                            // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
            Headers: amqp.Table{
                "x-error-type":    errorType,
                "x-error-message": errorMsg,
                "x-timestamp":     time.Now().Unix(),
            },
        },
    )
    
    if err != nil {
        log.Printf("Failed to publish to dead letter exchange: %v", err)
    } else {
        log.Printf("Message published to dead letter exchange with error type: %s", errorType)
    }
}

func (c *Consumer) Close() error {
    var err error

    if c.channel != nil {
        if err = c.channel.Close(); err != nil {
            log.Printf("Error closing channel: %v", err)
        }
    }

    if c.connection != nil {
        if err = c.connection.Close(); err != nil {
            log.Printf("Error closing connection: %v", err)
        }
    }

    return err
}