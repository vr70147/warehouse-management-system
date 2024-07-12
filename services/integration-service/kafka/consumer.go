package kafka

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func ConsumeMessages(topic string, handleMessage func(kafka.Message)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Topic:   topic,
		GroupID: "integration-service-group",
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("could not read message: %v", err)
			continue
		}

		log.Printf("Consumed message from topic %s: %s with account_id: %s", topic, string(m.Value), getAccountIDFromHeaders(m.Headers))
		handleMessage(m)
	}
}

// Helper function to get account_id from message headers
func getAccountIDFromHeaders(headers []kafka.Header) string {
	for _, h := range headers {
		if h.Key == "account_id" {
			return string(h.Value)
		}
	}
	return ""
}
