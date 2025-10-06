package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Channel *amqp091.Channel
	Queue   amqp091.Queue
}

func NewPublisher(r *RabbitMQ, queueName string) (*Publisher, error) {
	queue, err := r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		Channel: r.Channel,
		Queue:   queue,
	}, nil
}

func (p *Publisher) Publish(body []byte) error {
	return p.Channel.Publish(
		"",           // exchange
		p.Queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
