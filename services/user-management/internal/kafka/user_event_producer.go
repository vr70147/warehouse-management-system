package kafka

import (
	"encoding/json"
	"user-management/internal/model"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ProducerUserEvent(event model.UserEvent) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	eventBytes, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	topic := "user-events"

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          eventBytes,
	}, nil)

	p.Flush(15 * 1000)
}
