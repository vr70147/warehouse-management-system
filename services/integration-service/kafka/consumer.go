package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func ConsumeMessages(topic string, handleMessage func(kafka.Message)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
		GroupID: "integration-service-group",
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("could not read message: %v", err)
			continue
		}

		handleMessage(m)
	}
}
