package events

import "time"

// Envelope contains common metadata for every event published to Kafka.
//
// The event payload is serialized as JSON together with this metadata.
type Envelope[T any] struct {
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	OccurredAt time.Time `json:"occurred_at"`
	Producer   string    `json:"producer"`
	RunID      string    `json:"run_id"`
	Payload    T         `json:"payload"`
}
