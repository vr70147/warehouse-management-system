package main

import (
	"log"
	"user-management/internal/api"
	"user-management/internal/initializers"
	"user-management/internal/kafka"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabse()
}

func main() {
	api.Routers()

	broker := "localhost:9092"
	topic := "user-events"
	event := kafka.UserEvent{
		Type:     "UserCreated",
		UserID:   1,
		Username: "john_doe",
		Email:    "john@example.com",
		RoleID:   2,
	}

	if err := kafka.ProducerUserEvent(broker, topic, event); err != nil {
		log.Fatalf("Failed to produce user event: %v", err)
	}
	log.Println("Produced user event")
}
