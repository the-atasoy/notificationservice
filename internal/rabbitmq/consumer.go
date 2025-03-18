package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    queue   string
}

func NewConsumer(uri, queue string) (*Consumer, error) {
    conn, err := amqp.Dial(uri)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    return &Consumer{
        conn:    conn,
        channel: ch,
        queue:   queue,
    }, nil
}

func (c *Consumer) Start(handler func([]byte) error) error {
    msgs, err := c.channel.Consume(
        c.queue, // queue
        "",      // consumer
        true,    // auto-ack
        false,   // exclusive
        false,   // no-local
        false,   // no-wait
        nil,     // args
    )
    if err != nil {
        return err
    }

    go func() {
        for msg := range msgs {
            if err := handler(msg.Body); err != nil {
                log.Printf("Error processing message: %v", err)
            }
        }
    }()

    return nil
}