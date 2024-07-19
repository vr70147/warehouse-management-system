package events

import (
	"encoding/json"
	"log"
)

type Event struct {
	Data      map[string]interface{} `json:"data"`
	EventType string                 `json:"event_type"`
}

func (e *Event) ToJSON() []byte {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		log.Fatalf("failed to marshal event: %v", err)
	}

	return jsonBytes
}

func FromJSON(data []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
