package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func ProduceMessage() error {
	fmt.Printf("start producer")
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "test-topic",
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	return writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("Key-A"),
		Value: []byte("Hello from Go!"),
		Time:  time.Now(),
	})
}

func ConsumeMessage() error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test-topic",
		GroupID: "test-group",
	})
	defer reader.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		return err
	}

	log.Printf("Received message: key=%s value=%s\n", string(msg.Key), string(msg.Value))
	return nil
}

func Kafka_demo() {
	if err := ProduceMessage(); err != nil {
		log.Fatalf("Produce failed: %v", err)
	}
	if err := ConsumeMessage(); err != nil {
		log.Fatalf("Consume failed: %v", err)
	}
}
