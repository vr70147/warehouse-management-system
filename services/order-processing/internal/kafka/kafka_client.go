package kafka

import (
	"context"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/segmentio/kafka-go"
)

var (
	OrderWriter     *kafka.Writer
	InventoryWriter *kafka.Writer
	LowStockWriter  *kafka.Writer
)

// InitKafkaWriters initializes Kafka writers for order, inventory, and low stock notifications.
func InitKafkaWriters() {
	// Get broker addresses and topic names from environment variables.
	brokers := os.Getenv("KAFKA_BROKERS")
	orderTopic := os.Getenv("ORDER_EVENT_TOPIC")
	inventoryTopic := os.Getenv("INVENTORY_STATUS_TOPIC")
	lowStockTopic := os.Getenv("LOW_STOCK_NOTIFICATION_TOPIC")

	// Check if any of the required environment variables are not set.
	if brokers == "" || orderTopic == "" || inventoryTopic == "" || lowStockTopic == "" {
		log.Fatalf("KAFKA_BROKERS, ORDER_EVENT_TOPIC, INVENTORY_STATUS_TOPIC or LOW_STOCK_NOTIFICATION_TOPIC environment variable not set")
	}

	// Create a new Kafka admin client.
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0 // Adjust this based on your Kafka version

	admin, err := sarama.NewClusterAdmin([]string{brokers}, config)
	if err != nil {
		log.Fatalf("failed to create Kafka admin client: %v", err)
	}
	defer admin.Close()

	// Create topics if they do not exist.
	topics := []string{orderTopic, inventoryTopic, lowStockTopic}
	for _, topic := range topics {
		err = createTopicIfNotExists(admin, topic)
		if err != nil {
			log.Fatalf("failed to create topic %s: %v", topic, err)
		}
	}

	// Initialize Kafka writers for each topic.
	OrderWriter = createKafkaWriter(brokers, orderTopic)
	InventoryWriter = createKafkaWriter(brokers, inventoryTopic)
	LowStockWriter = createKafkaWriter(brokers, lowStockTopic)

	// Test connection and topic availability by sending a test message.
	testKafkaWriter(OrderWriter, "order events")
	testKafkaWriter(InventoryWriter, "inventory events")
	testKafkaWriter(LowStockWriter, "low stock notifications")
}

// createKafkaWriter creates and returns a Kafka writer for a given topic.
func createKafkaWriter(brokers, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// createTopicIfNotExists creates a Kafka topic if it does not already exist.
func createTopicIfNotExists(admin sarama.ClusterAdmin, topic string) error {
	topics, err := admin.ListTopics()
	if err != nil {
		return err
	}

	if _, ok := topics[topic]; ok {
		log.Printf("topic %s already exists", topic)
		return nil
	}

	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		return err
	}

	log.Printf("topic %s created", topic)
	return nil
}

// testKafkaWriter sends a test message to the Kafka writer to ensure it is properly initialized.
func testKafkaWriter(writer *kafka.Writer, description string) {
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Value: []byte("test message"),
	})
	if err != nil {
		log.Fatalf("failed to initialize Kafka writer for %s: %v", description, err)
	}
}
