package ports

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/google/uuid"
)

type UserReadRepository interface {
	FindByID(ctx context.Context, tenantID string, id uuid.UUID) (*entities.UserRead, error)
	Upsert(ctx context.Context, dto *entities.UserRead) error
}
