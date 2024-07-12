package kafka

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func ProduceMessage(topic, accountID, message string) error {
	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("KAFKA_BROKERS")),
		Topic:        topic,
		BatchTimeout: 10 * time.Millisecond,
	}

	// Include accountID in the message headers
	msg := kafka.Message{
		Value: []byte(message),
		Headers: []kafka.Header{
			{Key: "account_id", Value: []byte(accountID)},
		},
	}

	err := w.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	log.Printf("Produced message to topic %s: %s with account_id: %s", topic, message, accountID)
	return nil
}
