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

func ConsumerOrderEvents() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    os.Getenv("ORDER_EVENT_TOPIC"),
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

		var event model.OrderEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		if event.Action == "create" {
			processOrderCreation(event)
		} else if event.Action == "cancel" {
			processOrderCancellation(event)
		}
	}
}

func processOrderCreation(event model.OrderEvent) {
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		log.Printf("Database transaction error: %v\n", tx.Error)
		return
	}

	var stock model.Stock
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("product_id = ?", event.ProductID).First(&stock).Error; err != nil {
		log.Printf("Error finding stock: %v\n", err)
		tx.Rollback()
		return
	}

	log.Printf("Processing order creation for OrderID: %d, ProductID: %d, Requested Quantity: %d, Available Stock: %d\n", event.OrderID, event.ProductID, event.Quantity, stock.Quantity)

	if stock.Quantity >= event.Quantity {
		stock.Quantity -= event.Quantity
		if err := tx.Save(&stock).Error; err != nil {
			log.Printf("Error updating stock: %v\n", err)
			tx.Rollback()
		} else {
			log.Printf("Stock updated successfully: %+v\n", stock)
			if err := tx.Commit().Error; err != nil {
				log.Printf("Error committing transaction: %v\n", err)
				tx.Rollback()
			} else {
				publishInventoryStatus(event.OrderID, event.ProductID, event.Quantity, "Ready for Shipping")
				if stock.Quantity <= uint(stock.LowStockThreshold) {
					publishLowStockNotification(stock.ProductID, stock.Quantity, stock.LowStockThreshold)
				}
			}
		}
	} else {
		log.Printf("Not enough stock for product_id %d\n", event.ProductID)
		tx.Rollback()
		publishInventoryStatus(event.OrderID, event.ProductID, event.Quantity, "Out of Stock")
	}
}

func processOrderCancellation(event model.OrderEvent) {
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		log.Printf("Database transaction error: %v\n", tx.Error)
		return
	}

	var stock model.Stock
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("product_id = ?", event.ProductID).First(&stock).Error; err != nil {
		log.Printf("Error finding stock: %v\n", err)
		tx.Rollback()
		return
	}

	stock.Quantity += event.Quantity
	if err := tx.Save(&stock).Error; err != nil {
		log.Printf("Error updating stock: %v\n", err)
		tx.Rollback()
	} else {
		log.Printf("Stock updated successfully after cancellation: %+v\n", stock)
		if err := tx.Commit().Error; err != nil {
			log.Printf("Error committing transaction: %v\n", err)
			tx.Rollback()
		} else {
			publishInventoryStatus(event.OrderID, event.ProductID, event.Quantity, "Cancelled")
		}
	}
}

func publishInventoryStatus(orderID uint, productID uint, quantity uint, status string) {
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("INVENTORY_STATUS_TOPIC")

	if brokers == "" || topic == "" {
		log.Fatalf("KAFKA_BROKERS or INVENTORY_STATUS_TOPIC environment variable not set")
	}

	writer := kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	event := model.OrderEvent{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		Action:    status,
	}

	messageBytes, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("failed to marshal message: %v", err)
	}

	log.Printf("Sending message: %s\n", string(messageBytes))

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("order_id"),
		Value: messageBytes,
	})

	if err != nil {
		log.Fatalf("failed to write message to kafka: %v", err)
	}
}

func publishLowStockNotification(productID uint, quantity uint, lowStockThreshold int) {
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("LOW_STOCK_TOPIC")

	if brokers == "" || topic == "" {
		log.Fatalf("KAFKA_BROKERS or LOW_STOCK_TOPIC environment variable not set")
	}

	writer := kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	message := map[string]interface{}{
		"product_id":          productID,
		"quantity":            quantity,
		"low_stock_threshold": lowStockThreshold,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("failed to marshal message: %v", err)
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("product_id"),
		Value: messageBytes,
	})

	if err != nil {
		log.Fatalf("failed to write message to kafka: %v", err)
	}

	log.Printf("Low stock notification published for ProductID: %d\n", productID)
}
