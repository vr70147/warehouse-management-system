package eventbus

import (
	"context"
	"log"
	"warehouse-management-system/internal/events"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type EventBus struct {
	client *redis.Client
}

func NewEventBus(redisAddr string) *EventBus {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &EventBus{client: client}
}

func (eb *EventBus) Publish(channel string, event *events.Event) {
	err := eb.client.Publish(ctx, channel, event.ToJSON()).Err()
	if err != nil {
		log.Fatalf("Error publishing event: %v", err)
	}
}

func (eb *EventBus) Subscribe(channel string, handler func(event *events.Event)) {
	pubsub := eb.client.Subscribe(ctx, channel)
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			event, err := events.FromJSON([]byte(msg.Payload))
			if err != nil {
				log.Printf("Error unmarshalling event: %v", err)
				continue
			}
			handler(event)
		}
	}()
}
