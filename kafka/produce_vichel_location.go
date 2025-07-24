package kafka

import (
	model "FMTS/internal/tracking/domain/entity"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  brokers,
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
			// RequiredAcks: kafka.RequireAll,
		}),
	}
}

func (kp *KafkaProducer) ProduceVehicleLocation(ctx context.Context, location model.VehicleLocation) error {
	// Marshal the data to JSON
	value, err := json.Marshal(location)
	if err != nil {
		return err
	}

	// Create the Kafka message
	msg := kafka.Message{
		Key:   []byte(location.VehicleID), // useful for partitioning
		Value: value,
		Time:  time.Now(),
	}

	// Send message
	err = kp.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Error producing message to Kafka: %v", err)
		return err
	}

	log.Printf("Produced location for vehicle %s at %v", location.VehicleID, location.Timestamp)
	return nil
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}
