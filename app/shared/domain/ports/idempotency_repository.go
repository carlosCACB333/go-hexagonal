package ports

import "context"

type IdempotencyRepository interface {
	IsProcessed(ctx context.Context, tenantID, key string) (bool, error)
	MarkAsProcessed(ctx context.Context, tenantID, key string) error
}
