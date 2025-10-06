package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(cfg *KafkaConfig) (*KafkaProducer, error) {
	kafkaCfg := &kafka.ConfigMap{
		"bootstrap.servers":  cfg.BootstrapServers,
		"acks":               cfg.Acks,
		"enable.idempotence": cfg.EnableIdempotence,
	}

	if cfg.SecurityConfig != nil && cfg.SecurityConfig.EnableSASL {
		kafkaCfg.SetKey("security.protocol", "SASL_SSL")
		kafkaCfg.SetKey("sasl.mechanisms", cfg.SecurityConfig.SaslMechanism)
		kafkaCfg.SetKey("sasl.username", cfg.SecurityConfig.SaslUsername)
		kafkaCfg.SetKey("sasl.password", cfg.SecurityConfig.SaslPassword)
	}

	producer, err := kafka.NewProducer(kafkaCfg)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: producer}, nil
}

func (kp *KafkaProducer) Produce(topic string, key, value []byte) error {
	return kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}, nil)
}

func (kp *KafkaProducer) Close() {
	kp.producer.Flush(5000)
	kp.producer.Close()
}
