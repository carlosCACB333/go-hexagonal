package ports

import (
	"context"
)

type EventBus interface {
	Publish(ctx context.Context, event DomainEvent, correlationID string) error
	Close()
}
