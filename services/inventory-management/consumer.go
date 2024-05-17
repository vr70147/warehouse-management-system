package main

import (
	"inventory-management/internal/kafka"
	"log"
)

func Consumer() {
	broker := "localhost:9092"
	groupID := "inventory-management"
	topics := []string{"user-events"}

	kafka.ConsumeUserEvents(broker, groupID, topics)
	log.Println("Started consuming user events")
}
