package ports

import "time"

type DomainEvent interface {
	EventID() string
	EventType() string
	AggregateID() string
	OccurredOn() time.Time
}
