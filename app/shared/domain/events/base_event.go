package events

import (
	"time"

	"github.com/google/uuid"
)

type BaseEvent struct {
	EventIDValue     string    `json:"event_id"`
	EventTypeValue   string    `json:"event_type"`
	AggregateIDValue string    `json:"aggregate_id"`
	OccurredOnValue  time.Time `json:"occurred_on"`
}

func NewBaseEvent(eventType, aggregateID string) BaseEvent {
	return BaseEvent{
		EventIDValue:     uuid.New().String(),
		EventTypeValue:   eventType,
		AggregateIDValue: aggregateID,
		OccurredOnValue:  time.Now(),
	}
}

func (e BaseEvent) EventID() string {
	return e.EventIDValue
}

func (e BaseEvent) EventType() string {
	return e.EventTypeValue
}

func (e BaseEvent) AggregateID() string {
	return e.AggregateIDValue
}

func (e BaseEvent) OccurredOn() time.Time {
	return e.OccurredOnValue
}
