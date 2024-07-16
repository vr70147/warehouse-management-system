package kafka

import (
	"context"
	"encoding/json"
	"inventory-management/internal/initializers"
	"inventory-management/internal/model"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

// ConsumerOrderEvents consumes order events from Kafka and updates order status in the inventory system
func ConsumerOrderEvents() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    "order-events",
		GroupID:  "inventory-management-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			continue
		}
		log.Printf("Received message: %s\n", string(m.Value))

		var order model.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		if results := initializers.DB.First(&order, order.ID); results.Error != nil {
			log.Printf("Failed to find order: %v\n", results.Error)
			continue
		}

		if order.Quantity <= 10 {
			order.Status = "Ready"
		} else {
			order.Status = "Out of Stock"
		}

		if err := initializers.DB.Save(&order).Error; err != nil {
			log.Printf("Failed to update order status: %v\n", err)
			continue
		}

		if err := PublishOrderStatus(order.ID, order.Status); err != nil {
			log.Printf("Failed to publish order status: %v\n", err)
		}
	}
}

// PublishOrderStatus publishes the order status to Kafka
func PublishOrderStatus(orderID uint, status string) error {
	w := kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_BROKERS")),
		Topic:    "inventory-status",
		Balancer: &kafka.LeastBytes{},
	}

	message := map[string]interface{}{
		"orderID": orderID,
		"status":  status,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("orderID"),
		Value: messageBytes,
	})

	if err != nil {
		return err
	}

	return nil
}
