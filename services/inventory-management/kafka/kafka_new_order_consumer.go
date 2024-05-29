package kafka

import (
	"context"
	"inventory-management/internal/initializers"
	"inventory-management/internal/model"
	"log"

	"github.com/segmentio/kafka-go"
)

func ConsumerOrderEvents() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "order-events",
		GroupID:  "inventory-management-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			panic("could not read message " + err.Error())
		}
		log.Printf("received message: %s\n", string(m.Value))

		var order model.Order
		if results := initializers.DB.First(&order, string(m.Value)); results.Error != nil {
			log.Printf("failed to find order: %v", results.Error)
			continue
		}

		if order.Quantity <= 10 {
			order.Status = "Ready"
		} else {
			order.Status = "Out of Stock"
		}
		initializers.DB.Save(&order)

		PublishOrderStatus(order.ID, order.Status)
	}
}

func PublishOrderStatus(orderID uint, status string) {
	w := kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "inventory-status",
		Balancer: &kafka.LeastBytes{},
	}

	err := w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("orderID"),
		Value: []byte(status),
	})

	if err != nil {
		log.Fatalf("failed to write message to kafka: %v", err)
	}
}
