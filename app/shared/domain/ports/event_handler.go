package ports

import "context"

type EventHandler interface {
	HandleEvent(ctx context.Context, eventType string, data []byte) error
}
