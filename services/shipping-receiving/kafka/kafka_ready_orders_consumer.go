package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"shipping-receiving/internal/initializers"
	"shipping-receiving/internal/model"
	"shipping-receiving/internal/utils"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// ConsumeInventoryStatus consumes inventory status messages from Kafka and updates the shipping status accordingly
func ConsumeInventoryStatus() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    "inventory-status",
		GroupID:  "shipping-receiving-group",
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

		var statusMessage struct {
			OrderID uint   `json:"order_id"`
			Status  string `json:"status"`
		}

		if err := json.Unmarshal(m.Value, &statusMessage); err != nil {
			log.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		if err := updateShippingStatus(statusMessage.OrderID, statusMessage.Status); err != nil {
			log.Printf("Failed to update shipping status: %v\n", err)
			continue
		}
	}
}

// updateShippingStatus updates the shipping status based on the received inventory status
func updateShippingStatus(orderID uint, status string) error {
	var shipping model.Shipping
	if result := initializers.DB.First(&shipping, orderID); result.Error != nil {
		return result.Error
	}

	if status == "Ready" {
		shipping.Status = "Shipped"
	} else {
		shipping.Status = "Cannot Ship"
	}

	if result := initializers.DB.Save(&shipping); result.Error != nil {
		return result.Error
	}

	if shipping.Status == "Shipped" {
		emailSubject := "Your Order Has Been Shipped"
		emailBody := "Your order with Shipping ID " + strconv.Itoa(int(shipping.ID)) + " has been shipped successfully."
		if err := utils.SendEmail("customer@example.com", emailSubject, emailBody); err != nil {
			log.Printf("Failed to send notification email: %v\n", err)
			return err
		}
	}

	return publishShippingStatus(shipping.ID, shipping.Status)
}

// publishShippingStatus publishes the updated shipping status to Kafka
func publishShippingStatus(shippingID uint, status string) error {
	w := kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_BROKERS")),
		Topic:    "shipping-status",
		Balancer: &kafka.LeastBytes{},
	}

	message := struct {
		ShippingID uint   `json:"shipping_id"`
		Status     string `json:"status"`
	}{
		ShippingID: shippingID,
		Status:     status,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("shippingID"),
		Value: messageBytes,
	})

	if err != nil {
		return err
	}

	return nil
}
