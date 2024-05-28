package kafka

import (
	"context"
	"log"
	"order-processing/internal/initializers"
	"order-processing/internal/model"

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

func ConsumerOrderEvent() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "order-events",
		GroupID:  "order-processing-group",
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

		order.Status = "Shipped"
		initializers.DB.Save(&order)

		log.Printf("Order %s shipped successfully", string(m.Key))
	}
}
