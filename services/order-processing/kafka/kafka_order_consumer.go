package kafka

import (
	"context"
	"fmt"
	"log"
	"order-processing/internal/initializers"
	"order-processing/internal/model"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

const maxRetries = 3

func ConsumerOrderEvent() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    os.Getenv("ORDER_EVENTS_TOPIC"),
		GroupID:  "order-processing-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("could not read message: %v", err)
			continue
		}

		log.Printf("received message: %s\n", string(m.Value))

		var order model.Order
		if err := processOrderMessage(m.Value); err != nil {
			log.Printf("failed to process order: %v", err)
			continue
		}

		order.Status = "Shipped"
		initializers.DB.Save(&order)

		log.Printf("Order %s shipped successfully", string(m.Key))
	}
}

func processOrderMessage(msg []byte) error {
	var order model.Order
	for i := 0; i < maxRetries; i++ {
		if results := initializers.DB.First(&order, string(msg)); results.Error == nil {
			order.Status = "Shipped"
			initializers.DB.Save(&order)
			return nil
		} else {
			log.Printf("attempt %d failed to find order: %v", i+1, results.Error)
			time.Sleep(2 * time.Second)
		}
	}
	return fmt.Errorf("failed to process order after %d retries", maxRetries)
}
