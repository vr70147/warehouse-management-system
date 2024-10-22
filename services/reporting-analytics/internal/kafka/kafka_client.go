package kafka

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

var (
	SalesWriter        *kafka.Writer
	InventoryWriter    *kafka.Writer
	ShippingWriter     *kafka.Writer
	UserActivityWriter *kafka.Writer
)

// InitKafkaWriters initializes Kafka writers for different topics
func InitKafkaWriters() {
	brokers := os.Getenv("KAFKA_BROKERS")
	salesTopic := os.Getenv("SALES_EVENT_TOPIC")
	inventoryTopic := os.Getenv("INVENTORY_LEVEL_TOPIC")
	shippingTopic := os.Getenv("SHIPPING_STATUS_TOPIC")
	userActivityTopic := os.Getenv("USER_ACTIVITY_TOPIC")

	if brokers == "" || salesTopic == "" || inventoryTopic == "" || shippingTopic == "" || userActivityTopic == "" {
		log.Fatalf("KAFKA_BROKERS or one of the topic environment variables is not set")
	}

	// Create topics if they don't exist
	createTopic(brokers, salesTopic)
	createTopic(brokers, inventoryTopic)
	createTopic(brokers, shippingTopic)
	createTopic(brokers, userActivityTopic)

	// Initialize Kafka writer for sales events
	SalesWriter = createKafkaWriter(brokers, salesTopic)
	// Initialize Kafka writer for inventory levels
	InventoryWriter = createKafkaWriter(brokers, inventoryTopic)
	// Initialize Kafka writer for shipping statuses
	ShippingWriter = createKafkaWriter(brokers, shippingTopic)
	// Initialize Kafka writer for user activities
	UserActivityWriter = createKafkaWriter(brokers, userActivityTopic)

	// Test the Kafka writers by sending a test message
	testKafkaWriter(SalesWriter, "sales events")
	testKafkaWriter(InventoryWriter, "inventory levels")
	testKafkaWriter(ShippingWriter, "shipping statuses")
	testKafkaWriter(UserActivityWriter, "user activities")
}

// createKafkaWriter creates and returns a Kafka writer for a given topic
func createKafkaWriter(brokers, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// testKafkaWriter sends a test message to ensure the Kafka writer is properly initialized
func testKafkaWriter(writer *kafka.Writer, description string) {
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Value: []byte("test message"),
	})
	if err != nil {
		log.Fatalf("failed to initialize Kafka writer for %s: %v", description, err)
	}
}

// createTopic creates a Kafka topic if it doesn't exist
func createTopic(brokers, topic string) {
	conn, err := kafka.Dial("tcp", brokers)
	if err != nil {
		log.Fatalf("failed to dial kafka: %v", err)
	}
	defer conn.Close()

	partitionCount := 1
	replicationFactor := 1

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitionCount,
		ReplicationFactor: replicationFactor,
	})
	if err != nil && err != kafka.TopicAlreadyExists {
		log.Fatalf("failed to create topic %s: %v", topic, err)
	}
}
