package config

import (
	"log"
	"os"
	"sync"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaWriter *kafka.Writer
	kafkaOnce   sync.Once
)

// GetKafkaWriter returns a singleton Kafka writer (producer)
func GetKafkaWriter() *kafka.Writer {
	kafkaOnce.Do(func() {
		broker := os.Getenv("KAFKA_BROKER")
		topic := os.Getenv("KAFKA_TOPIC")
		if broker == "" {
			broker = "localhost:9092"
		}
		if topic == "" {
			topic = "elastic"
		}
		kafkaWriter = &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}
		log.Println("Connected to Kafka")
	})
	return kafkaWriter
}
