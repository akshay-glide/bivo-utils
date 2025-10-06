package rabbitmq

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler func(d amqp091.Delivery)

type Consumer struct {
	Channel  *amqp091.Channel
	Queue    amqp091.Queue
	Consumer string
}

func NewConsumer(r *RabbitMQ, queueName, consumerTag string) (*Consumer, error) {
	queue, err := r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Channel:  r.Channel,
		Queue:    queue,
		Consumer: consumerTag,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) error {
	msgs, err := c.Channel.Consume(
		c.Queue.Name,
		c.Consumer,
		false, // auto-ack: false for manual
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					log.Println("RabbitMQ consumer channel closed")
					return
				}
				handler(d)
			case <-ctx.Done():
				log.Println("Consumer context canceled")
				return
			}
		}
	}()

	return nil
}
