package rabbitmq

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestRabbitMQPubSub(t *testing.T) {
	// Context for test timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queueName := "test-queue"
	testMessage := []byte("hello-rabbitmq")

	// Connect
	conn, err := NewConnection("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Declare publisher
	publisher, err := NewPublisher(conn, queueName)
	if err != nil {
		t.Fatalf("Failed to create publisher: %v", err)
	}

	// Declare consumer
	consumer, err := NewConsumer(conn, queueName, "test-consumer")
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	// Channel to signal message received
	received := make(chan []byte, 1)

	// Start consuming
	err = consumer.Consume(ctx, func(msg amqp091.Delivery) {
		log.Printf("Received message: %s", msg.Body)
		received <- msg.Body
		_ = msg.Ack(false)
	})
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	// Give RabbitMQ a moment to register consumer
	time.Sleep(500 * time.Millisecond)

	// Publish message
	err = publisher.Publish(testMessage)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Wait for message or timeout
	select {
	case msg := <-received:
		if string(msg) != string(testMessage) {
			t.Errorf("Expected message %s, got %s", testMessage, msg)
		}
	case <-ctx.Done():
		t.Fatal("Timeout waiting for RabbitMQ message")
	}
}
