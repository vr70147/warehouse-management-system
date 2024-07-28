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
		Topic:    "sales-events",
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
