package kafka

// import (
// 	model "FMTS/internal/tracking/domain/entity"
// 	// "FMTS/internal/tracking/domain/entity"
// 	// "FMTS/internal/tracking/infrastructure/persistence"
// 	"context"
// 	"encoding/json"
// 	"log"

// 	"github.com/segmentio/kafka-go"
// )

// type KafkaConsumer struct {
// 	Reader     *kafka.Reader
// 	Repository *persistence.VehicleLocationRepository
// }

// func NewKafkaConsumer(brokers []string, topic, groupID string, repo *persistence.VehicleLocationRepository) *KafkaConsumer {
// 	reader := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:     brokers,
// 		Topic:       topic,
// 		GroupID:     groupID,
// 		StartOffset: kafka.LastOffset,
// 		MinBytes:    10e3, // 10KB
// 		MaxBytes:    10e6, // 10MB
// 	})
// 	return &KafkaConsumer{
// 		Reader:     reader,
// 		Repository: repo,
// 	}
// }

// func (kc *KafkaConsumer) StartConsuming(ctx context.Context) {
// 	for {
// 		m, err := kc.Reader.ReadMessage(ctx)
// 		if err != nil {
// 			log.Printf("Error reading message: %v", err)
// 			continue
// 		}

// 		var location model.VehicleLocation
// 		if err := json.Unmarshal(m.Value, &location); err != nil {
// 			log.Printf("Error unmarshalling message: %v", err)
// 			continue
// 		}

// 		if err := kc.Repository.Save(ctx, &location); err != nil {
// 			log.Printf("Error saving to DB: %v", err)
// 		}
// 	}
// }
