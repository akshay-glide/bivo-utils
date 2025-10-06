package kafka

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var testKafkaConfig = &KafkaConfig{
	BootstrapServers:  "localhost:9092",
	GroupID:           "user-service-group",
	Topic:             "user-created-topic",
	AutoOffset:        "earliest",
	Acks:              "all",
	EnableIdempotence: true,
}

func TestKafkaProduceAndConsume(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start consumer
	consumer, err := NewKafkaConsumer(testKafkaConfig)
	if err != nil {
		t.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	err = consumer.Subscribe()
	if err != nil {
		t.Fatalf("Failed to subscribe to topic: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	received := make(chan *kafka.Message, 1)

	err = consumer.Start(ctx, func(msg *kafka.Message) {
		log.Printf("Received Kafka message: key=%s value=%s", msg.Key, msg.Value)
		received <- msg
		wg.Done()
	})
	if err != nil {
		t.Fatalf("Failed to start Kafka consumer: %v", err)
	}

	// Wait a bit to ensure consumer is ready
	time.Sleep(5 * time.Second)

	// Produce message
	producer, err := NewKafkaProducer(testKafkaConfig)
	if err != nil {
		t.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	key := []byte("test-new-msg")
	value := []byte("hello-kafka")

	err = producer.Produce(testKafkaConfig.Topic, key, value)
	if err != nil {
		t.Fatalf("Failed to produce Kafka message: %v", err)
	}

	// Wait for message or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		msg := <-received
		if string(msg.Key) != string(key) {
			t.Errorf("Expected key %s, got %s", key, msg.Key)
		}
		if string(msg.Value) != string(value) {
			t.Errorf("Expected value %s, got %s", value, msg.Value)
		}
		log.Println("Kafka produce and consume test passed")
	case <-ctx.Done():
		t.Fatal("Timeout waiting for Kafka message")
	}
}
