package message_broker

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func PublishOrderEvent(orderID string) {
	w := kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_BROKERS")),
		Topic:    os.Getenv("ORDER_EVENTS_TOPIC"),
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
