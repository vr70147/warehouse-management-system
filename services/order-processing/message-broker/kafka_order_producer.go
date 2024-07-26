package message_broker

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func PublishOrderEvent(orderID uint, status string) {
	w := kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_BROKERS")),
		Topic:    os.Getenv("ORDER_EVENT_TOPIC"),
		Balancer: &kafka.LeastBytes{},
	}

	message := map[string]interface{}{
		"order_id": orderID,
		"status":   status,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("failed to marshal message: %v", err)
	}

	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("order_id"),
		Value: messageBytes,
	})

	if err != nil {
		log.Fatalf("failed to write message to kafka: %v", err)
	}
}
