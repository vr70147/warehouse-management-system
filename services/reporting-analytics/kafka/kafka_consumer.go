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
			log.Fatal("could not read message " + err.Error())
		}
		log.Printf("received message: %s\n", string(m.Value))

		var salesReport model.SalesReport
		err = json.Unmarshal(m.Value, &salesReport)

		if err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		salesReport.Timestamp = time.Now()

		if results := initializers.DB.Create(&salesReport); results.Error != nil {
			log.Printf("failed to create sales report: %v", results.Error)
			continue
		}

		log.Printf("Sales report created successfully: %v", salesReport)
	}
}
