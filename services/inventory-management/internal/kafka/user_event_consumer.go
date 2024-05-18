package kafka

import (
	"encoding/json"
	"inventory-management/internal/model"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumeUserEvents(broker string, groupID string, topics []string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"groupID":           groupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer c.Close()

	c.SubscribeTopics(topics, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			var event model.UserEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Failed to unmarshal event: %s", err)
				continue
			}
			handleUserEvent(event)
		}
	}
}

func handleUserEvent(event model.UserEvent) {
	switch event.EventType {
	case "UserCreated":
		log.Printf("Handling user created event: %+v\n", event)
	case "UserUpdated":
		log.Printf("Handling user updated event: %+v\n", event)
	default:
		log.Printf("Unhandled event type: %s\n", event.EventType)
	}
}
