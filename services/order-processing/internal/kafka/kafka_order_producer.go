package kafka

import (
	"context"
	"encoding/json"
	"log"
	"order-processing/internal/model"

	"github.com/segmentio/kafka-go"
)

// PublishOrderEvent publishes an order event to the Kafka topic.
func PublishOrderEvent(orderID uint, productID uint, quantity uint, action string) {
	// Create an order event struct
	event := model.OrderEvent{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		Action:    action,
	}

	// Marshal the order event into JSON
	messageBytes, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("failed to marshal order event: %v", err)
	}

	// Write the JSON message to the Kafka topic
	err = OrderWriter.WriteMessages(context.Background(), kafka.Message{
		Value: messageBytes,
	})
	if err != nil {
		log.Fatalf("failed to write message to kafka: %v", err)
	}

	log.Printf("Published order event: %+v\n", event)
}
