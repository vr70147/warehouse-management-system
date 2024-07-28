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

func InitKafkaWriters() {
	brokers := os.Getenv("KAFKA_BROKERS")
	orderTopic := os.Getenv("ORDER_EVENT_TOPIC")
	inventoryTopic := os.Getenv("INVENTORY_STATUS_TOPIC")
	lowStockTopic := os.Getenv("LOW_STOCK_NOTIFICATION_TOPIC")

	if brokers == "" || orderTopic == "" || inventoryTopic == "" || lowStockTopic == "" {
		log.Fatalf("KAFKA_BROKERS, ORDER_EVENT_TOPIC, INVENTORY_STATUS_TOPIC or LOW_STOCK_NOTIFICATION_TOPIC environment variable not set")
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0 // Adjust this based on your Kafka version

	admin, err := sarama.NewClusterAdmin([]string{brokers}, config)
	if err != nil {
		log.Fatalf("failed to create Kafka admin client: %v", err)
	}
	defer admin.Close()

	topics := []string{orderTopic, inventoryTopic, lowStockTopic}
	for _, topic := range topics {
		err = createTopicIfNotExists(admin, topic)
		if err != nil {
			log.Fatalf("failed to create topic %s: %v", topic, err)
		}
	}

	OrderWriter = &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    orderTopic,
		Balancer: &kafka.LeastBytes{},
	}

	InventoryWriter = &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    inventoryTopic,
		Balancer: &kafka.LeastBytes{},
	}

	LowStockWriter = &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    lowStockTopic,
		Balancer: &kafka.LeastBytes{},
	}

	// Test connection and topic availability
	err = OrderWriter.WriteMessages(context.Background(), kafka.Message{
		Value: []byte("test message"),
	})
	if err != nil {
		log.Fatalf("failed to initialize Kafka writer for order events: %v", err)
	}

	err = InventoryWriter.WriteMessages(context.Background(), kafka.Message{
		Value: []byte("test message"),
	})
	if err != nil {
		log.Fatalf("failed to initialize Kafka writer for inventory events: %v", err)
	}

	err = LowStockWriter.WriteMessages(context.Background(), kafka.Message{
		Value: []byte("test message"),
	})
	if err != nil {
		log.Fatalf("failed to initialize Kafka writer for low stock notifications: %v", err)
	}
}

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
