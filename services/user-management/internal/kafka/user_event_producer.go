package kafka

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type UserEvent struct {
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   uint   `json:"role_id"`
}

func ProducerUserEvent(broker string, topic string, event UserEvent) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return err
	}
	defer p.Close()

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          eventBytes,
	}, nil)

	p.Flush(15 * 1000)
	return nil
}
