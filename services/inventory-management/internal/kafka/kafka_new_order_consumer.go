package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"inventory-management/internal/initializers"
	"inventory-management/internal/model"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

var ErrEmptyQueue = errors.New("kafka: empty queue")

type KafkaService interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type KafkaClient struct {
	Reader *kafka.Reader
	Writer *kafka.Writer
}

func (kc *KafkaClient) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return kc.Reader.ReadMessage(ctx)
}

func (kc *KafkaClient) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	return kc.Writer.WriteMessages(ctx, msgs...)
}

func NewKafkaClient() *KafkaClient {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	fmt.Println(kafkaBrokers)
	newOrderEventsTopic := os.Getenv("NEW_ORDER_EVENTS_TOPIC")

	if kafkaBrokers == "" || newOrderEventsTopic == "" {
		log.Fatalf("Environment variables KAFKA_BROKERS and NEW_ORDER_EVENTS_TOPIC must be set")
	}

	log.Printf("Initializing Kafka client with brokers: %s and topic: %s", kafkaBrokers, newOrderEventsTopic)

	return &KafkaClient{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{kafkaBrokers},
			Topic:    newOrderEventsTopic,
			GroupID:  "inventory-management-group",
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(kafkaBrokers),
			Topic:    "inventory-status",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

var KafkaSvc KafkaService = NewKafkaClient()

func ConsumerOrderEvents() {
	log.Println("Starting Kafka consumer for order events")
	for {
		m, err := KafkaSvc.ReadMessage(context.Background())
		if err != nil {
			if err == ErrEmptyQueue {
				log.Println("No messages in the queue")
				continue
			}
			log.Printf("Error reading message: %v\n", err)
			continue
		}
		log.Printf("Received message: %s\n", string(m.Value))

		var order model.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		log.Printf("Processing order ID: %d", order.ID)

		if results := initializers.DB.First(&order, order.ID); results.Error != nil {
			log.Printf("Failed to find order: %v\n", results.Error)
			continue
		}

		log.Printf("Found order in DB: %+v", order)

		if order.Quantity <= 10 {
			order.Status = "Ready"
		} else {
			order.Status = "Out of Stock"
		}

		log.Printf("Updating order status to: %s", order.Status)

		if err := initializers.DB.Save(&order).Error; err != nil {
			log.Printf("Failed to update order status: %v\n", err)
			continue
		}

		log.Printf("Order status updated: %+v", order)

		if err := PublishOrderStatus(order.ID, order.Status); err != nil {
			log.Printf("Failed to publish order status: %v\n", err)
		} else {
			log.Printf("Order status published: OrderID=%d, Status=%s", order.ID, order.Status)
		}
	}
}

func PublishOrderStatus(orderID uint, status string) error {
	message := map[string]interface{}{
		"orderID": orderID,
		"status":  status,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	log.Printf("Publishing message to Kafka: %s", string(messageBytes))

	err = KafkaSvc.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("orderID"),
		Value: messageBytes,
	})

	if err != nil {
		return err
	}

	return nil
}
