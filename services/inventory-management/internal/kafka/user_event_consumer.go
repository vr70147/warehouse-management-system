package kafka

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type UserEvent struct {
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   uint   `json:"role_id"`
}

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
			var event UserEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Failed to unmarshal event: %s", err)
				continue
			}
			handleUserEvent(event)
		}
	}
}

func handleUserEvent(event UserEvent) {
	switch event.Type {
	case "UserCreated":
		log.Printf("Handling user created event: %+v\n", event)
	case "UserUpdated":
		log.Printf("Handling user updated event: %+v\n", event)
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}
}
