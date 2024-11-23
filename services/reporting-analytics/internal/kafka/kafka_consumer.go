package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"reporting-analytics/internal/initializers"
	"reporting-analytics/internal/model"
	"time"

	"github.com/segmentio/kafka-go"
)

// ConsumerSalesEvent consumes sales events from Kafka and processes them
func ConsumerSalesEvent() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    os.Getenv("SALES_EVENT_TOPIC"),
		GroupID:  "reporting-analytics-group",
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

		var salesReport model.SalesReport
		if err := json.Unmarshal(m.Value, &salesReport); err != nil {
			log.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		salesReport.Timestamp = time.Now()

		if err := saveSalesReport(&salesReport); err != nil {
			log.Printf("Failed to save sales report: %v\n", err)
			continue
		}

		log.Printf("Sales report created successfully: %v\n", salesReport)
	}
}

// saveSalesReport saves the sales report to the database
func saveSalesReport(salesReport *model.SalesReport) error {
	if result := initializers.DB.Create(salesReport); result.Error != nil {
		return result.Error
	}
	return nil
}

// ConsumerInventoryLevel consumes inventory level updates from Kafka and processes them
func ConsumerInventoryLevel() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    os.Getenv("INVENTORY_LEVEL_TOPIC"),
		GroupID:  "reporting-analytics-group",
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

		var inventoryLevel model.InventoryLevel
		if err := json.Unmarshal(m.Value, &inventoryLevel); err != nil {
			log.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		inventoryLevel.LastUpdated = time.Now()

		if err := saveInventoryLevel(&inventoryLevel); err != nil {
			log.Printf("Failed to save inventory level: %v\n", err)
			continue
		}

		log.Printf("Inventory level updated successfully: %v\n", inventoryLevel)
	}
}

// saveInventoryLevel saves the inventory level to the database
func saveInventoryLevel(inventoryLevel *model.InventoryLevel) error {
	if result := initializers.DB.Create(inventoryLevel); result.Error != nil {
		return result.Error
	}
	return nil
}

// ConsumerShippingStatus consumes shipping status updates from Kafka and processes them
func ConsumerShippingStatus() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    os.Getenv("SHIPPING_STATUS_TOPIC"),
		GroupID:  "reporting-analytics-group",
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

		var shippingStatus model.ShippingStatus
		if err := json.Unmarshal(m.Value, &shippingStatus); err != nil {
			log.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		shippingStatus.LastUpdated = time.Now()

		if err := saveShippingStatus(&shippingStatus); err != nil {
			log.Printf("Failed to save shipping status: %v\n", err)
			continue
		}

		log.Printf("Shipping status updated successfully: %v\n", shippingStatus)
	}
}

// saveShippingStatus saves the shipping status to the database
func saveShippingStatus(shippingStatus *model.ShippingStatus) error {
	if result := initializers.DB.Create(shippingStatus); result.Error != nil {
		return result.Error
	}
	return nil
}

// ConsumerUserActivity consumes user activity events from Kafka and processes them
func ConsumerUserActivity() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    os.Getenv("USER_ACTIVITY_TOPIC"),
		GroupID:  "reporting-analytics-group",
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

		var userActivity model.UserActivity
		if err := json.Unmarshal(m.Value, &userActivity); err != nil {
			log.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		userActivity.Timestamp = time.Now()

		if err := saveUserActivity(&userActivity); err != nil {
			log.Printf("Failed to save user activity: %v\n", err)
			continue
		}

		log.Printf("User activity recorded successfully: %v\n", userActivity)
	}
}

// saveUserActivity saves the user activity to the database
func saveUserActivity(userActivity *model.UserActivity) error {
	if result := initializers.DB.Create(userActivity); result.Error != nil {
		return result.Error
	}
	return nil
}
