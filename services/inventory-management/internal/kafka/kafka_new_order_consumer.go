package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"inventory-management/internal/model"
	"log"
	"net/http"
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
	topic := os.Getenv("ORDER_EVENT_TOPIC")
	if topic == "" {
		log.Fatalf("ORDER_EVENT_TOPIC environment variable not set")
	}
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatalf("KAFKA_BROKERS environment variable not set")
	}

	log.Printf("Initializing Kafka client with brokers: %s and topic: %s", kafkaBrokers, topic)

	return &KafkaClient{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{kafkaBrokers},
			Topic:    topic,
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

		var orderEvent struct {
			OrderID uint   `json:"order_id"`
			Status  string `json:"status"`
		}

		if err := json.Unmarshal(m.Value, &orderEvent); err != nil {
			log.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		log.Printf("Processing order ID: %d", orderEvent.OrderID)

		orderDetails, err := fetchOrderDetails(orderEvent.OrderID)
		if err != nil {
			log.Printf("Failed to fetch order details: %v\n", err)
			continue
		}

		log.Printf("Fetched order details: %+v", orderDetails)

		if orderDetails.Quantity <= 10 {
			orderDetails.Status = "Ready"
		} else {
			orderDetails.Status = "Out of Stock"
		}

		if err := updateOrderStatus(orderDetails.ID, orderDetails.Status); err != nil {
			log.Printf("Failed to update order status: %v\n", err)
			continue
		}

		log.Printf("Order status updated: %+v", orderDetails)

		if err := PublishOrderStatus(orderDetails.ID, orderDetails.Status); err != nil {
			log.Printf("Failed to publish order status: %v\n", err)
		} else {
			log.Printf("Order status published: OrderID=%d, Status=%s", orderDetails.ID, orderDetails.Status)
		}
	}
}

func fetchOrderDetails(orderID uint) (model.Order, error) {
	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	url := fmt.Sprintf("%s/orders/%d", orderServiceURL, orderID)

	resp, err := http.Get(url)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to fetch order details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Order{}, fmt.Errorf("failed to fetch order details: received status code %d", resp.StatusCode)
	}

	var order model.Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order details: %w", err)
	}

	return order, nil
}

func updateOrderStatus(orderID uint, status string) error {
	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	url := fmt.Sprintf("%s/orders/%d", orderServiceURL, orderID)

	orderUpdate := map[string]string{
		"status": status,
	}

	orderUpdateBytes, err := json.Marshal(orderUpdate)
	if err != nil {
		return fmt.Errorf("failed to marshal order update: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(orderUpdateBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update order status: received status code %d", resp.StatusCode)
	}

	return nil
}

func PublishOrderStatus(orderID uint, status string) error {
	message := map[string]interface{}{
		"order_id": orderID,
		"status":   status,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	log.Printf("Publishing message to Kafka: %s", string(messageBytes))

	err = KafkaSvc.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("order_id"),
		Value: messageBytes,
	})

	if err != nil {
		return err
	}

	return nil
}
