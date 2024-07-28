package kafka

import (
	"context"
	"encoding/json"
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
	topic := os.Getenv("ORDER_EVENT_TOPIC")
	if topic == "" {
		log.Fatalf("ORDER_EVENT_TOPIC environment variable not set")
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    topic,
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

		log.Printf("Order %s processed successfully", string(m.Key))
	}
}

func processOrderMessage(msg []byte) error {
	var orderEvent struct {
		OrderID   uint   `json:"order_id"`
		ProductID uint   `json:"product_id"`
		Quantity  int    `json:"quantity"`
		Action    string `json:"action"`
	}

	log.Printf("Processing order message: %s\n", string(msg))

	if err := json.Unmarshal(msg, &orderEvent); err != nil {
		return fmt.Errorf("failed to unmarshal order event: %v", err)
	}

	log.Printf("Unmarshalled order event: %+v\n", orderEvent)

	var order model.Order
	for i := 0; i < maxRetries; i++ {
		if results := initializers.DB.First(&order, orderEvent.OrderID); results.Error == nil {
			order.Status = orderEvent.Action
			if saveErr := initializers.DB.Save(&order).Error; saveErr != nil {
				return fmt.Errorf("failed to save order status: %v", saveErr)
			}
			return nil
		} else {
			log.Printf("attempt %d failed to find order: %v", i+1, results.Error)
			time.Sleep(2 * time.Second)
		}
	}
	return fmt.Errorf("failed to process order after %d retries", maxRetries)
}

func ConsumerInventoryStatus() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    os.Getenv("INVENTORY_STATUS_TOPIC"),
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

		var event model.OrderEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("failed to unmarshal inventory status: %v\n", err)
			continue
		}

		log.Printf("Unmarshalled Inventory Status: %+v\n", event)

		var order model.Order
		if result := initializers.DB.First(&order, event.OrderID); result.Error == nil {
			log.Printf("Updating order status for OrderID: %d, Status: %s\n", event.OrderID, event.Action)
			order.Status = event.Action
			if err := initializers.DB.Save(&order).Error; err != nil {
				log.Printf("Error updating order status: %v\n", err)
			} else {
				log.Printf("Order status updated successfully for OrderID: %d\n", event.OrderID)
				if event.Action == "Ready for Shipping" {
					PublishOrderEvent(order.ID, order.ProductID, order.Quantity, "ship")
				}
			}
		} else {
			log.Printf("Order not found: %v\n", result.Error)
		}
	}
}

func ConsumerShippingStatus() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    os.Getenv("SHIPPING_STATUS_TOPIC"),
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

		var shippingStatus struct {
			OrderID uint   `json:"order_id"`
			Action  string `json:"action"`
		}

		if err := json.Unmarshal(m.Value, &shippingStatus); err != nil {
			log.Printf("failed to unmarshal shipping status: %v", err)
			continue
		}

		var order model.Order
		if result := initializers.DB.First(&order, shippingStatus.OrderID); result.Error == nil {
			log.Printf("Updating order status for OrderID: %d, Action: %s\n", shippingStatus.OrderID, shippingStatus.Action)
			order.Status = shippingStatus.Action
			if err := initializers.DB.Save(&order).Error; err != nil {
				log.Printf("Error updating order status: %v\n", err)
			}
		}
	}
}
