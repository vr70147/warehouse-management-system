package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func PublishOrderEvent(orderID string) {
	w := kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "order-events",
		Balancer: &kafka.LeastBytes{},
	}
	err := w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("orderID"),
		Value: []byte(orderID),
	})

	if err != nil {
		log.Fatalf("failed to write message to kafka: %v", err)
	}
}
