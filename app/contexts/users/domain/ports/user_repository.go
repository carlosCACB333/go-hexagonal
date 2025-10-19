package ports

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/value_objects"
	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, tenantID string, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, tenantID string, email value_objects.Email) (*entities.User, error)
	ExistsByEmail(ctx context.Context, tenantID string, email value_objects.Email) (bool, error)
}
