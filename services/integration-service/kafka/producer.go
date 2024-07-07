package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func ProduceMessage(topic string, message string) error {
	w := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        topic,
		BatchTimeout: 10 * time.Millisecond,
	}

	err := w.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})

	if err != nil {
		return err
	}

	log.Printf("Producer message to topic %s: %s", topic, message)
	return nil
}
