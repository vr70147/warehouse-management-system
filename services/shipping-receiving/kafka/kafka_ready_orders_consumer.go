package kafka

import (
	"context"
	"log"
	"shipping-receiving/internal/initializers"
	"shipping-receiving/internal/model"

	"github.com/segmentio/kafka-go"
)

func ConsumeInventoryStatus() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "inventory-status",
		GroupID:  "shipping-receiving-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			panic("could not read message " + err.Error())
		}
		log.Printf("received message: %s\n", string(m.Value))

		var shipping model.Shipping
		if result := initializers.DB.First(&shipping, string(m.Key)); result.Error != nil {
			log.Printf("failed to find shipping: %v", result.Error)
			continue
		}

		if string(m.Value) == "Ready" {
			shipping.Status = "Shipped"
		} else {
			shipping.Status = "Cannot Ship"
		}

		initializers.DB.Save(&shipping)

		PublishShippingStatus(shipping.ID, shipping.Status)
	}
}

func PublishShippingStatus(shippingID uint, status string) {
	w := kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "shipping-status",
		Balancer: &kafka.LeastBytes{},
	}
	err := w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("shippingID"),
		Value: []byte(status),
	})

	if err != nil {
		log.Fatalf("failed to write message to kafka: %v", err)
	}
}
