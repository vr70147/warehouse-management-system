package kafka

import (
	"context"
	"encoding/json"
	"log"
	"order-processing/internal/model"

	"github.com/segmentio/kafka-go"
)

func PublishOrderEvent(orderID uint, productID uint, quantity uint, action string) {
	event := model.OrderEvent{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		Action:    action,
	}

	messageBytes, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("failed to marshal order event: %v", err)
	}

	err = OrderWriter.WriteMessages(context.Background(), kafka.Message{
		Value: messageBytes,
	})
	if err != nil {
		log.Fatalf("failed to write message to kafka: %v", err)
	}

	log.Printf("Published order event: %+v\n", event)
}
