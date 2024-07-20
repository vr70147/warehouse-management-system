package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"order-processing/internal/initializers"
	"order-processing/internal/model"
	"order-processing/internal/utils"
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

		if err := processOrderMessage(m.Value); err != nil {
			log.Printf("failed to process order: %v", err)
			continue
		}

		log.Printf("Order %s shipped successfully", string(m.Key))
	}
}

func processOrderMessage(msg []byte) error {
	var orderEvent struct {
		OrderID uint   `json:"order_id"`
		Status  string `json:"status"`
	}

	if err := json.Unmarshal(msg, &orderEvent); err != nil {
		return fmt.Errorf("failed to unmarshal order event: %v", err)
	}

	var order model.Order
	for i := 0; i < maxRetries; i++ {
		if results := initializers.DB.First(&order, orderEvent.OrderID); results.Error == nil {
			order.Status = orderEvent.Status
			initializers.DB.Save(&order)

			if orderEvent.Status == "Cancelled" {
				emailSubject := "Your Order Has Been Cancelled"
				emailBody := fmt.Sprintf("Your order with Order ID %d has been cancelled.", orderEvent.OrderID)
				if err := utils.SendEmail("customer@example.com", emailSubject, emailBody); err != nil {
					log.Printf("Failed to send notification email: %v\n", err)
					return err
				}
			}
			return nil
		} else {
			log.Printf("attempt %d failed to find order: %v", i+1, results.Error)
			time.Sleep(2 * time.Second)
		}
	}
	return fmt.Errorf("failed to process order after %d retries", maxRetries)
}
