package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"shipping-receiving/internal/initializers"
	"shipping-receiving/internal/model"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func ConsumerShippingEvents() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    os.Getenv("SHIPPING_EVENT_TOPIC"),
		GroupID:  "shipping-management-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("could not read message: %v", err)
			continue
		}

		var event struct {
			OrderID   uint   `json:"order_id"`
			ProductID uint   `json:"product_id"`
			Quantity  uint   `json:"quantity"`
			Action    string `json:"action"`
		}
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("failed to unmarshal shipping event: %v", err)
			continue
		}

		if event.Action == "ship" {
			var shipping model.Shipping
			shipping.OrderID = event.OrderID
			shipping.Status = "Shipped"
			shipping.ShippingDate = time.Now()

			if err := initializers.DB.Create(&shipping).Error; err != nil {
				log.Printf("failed to create shipping record: %v", err)
			}

			// Update the order status via HTTP request
			if err := updateOrderStatus(event.OrderID, "Shipped"); err != nil {
				log.Printf("failed to update order status: %v", err)
			}
		}
	}
}

func updateOrderStatus(orderID uint, status string) error {
	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	url := orderServiceURL + "/orders/" + strconv.Itoa(int(orderID)) + "/status"

	payload := map[string]string{
		"status": status,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update order status, status code: %d", resp.StatusCode)
	}

	return nil
}
