package queries

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/ports"
	"github.com/google/uuid"
)

type GetUserQuery struct {
	TenantID string
	UserID   uuid.UUID
}

type GetUserUseCase struct {
	readRepo ports.UserReadRepository
}

func NewGetUserUseCase(readModel ports.UserReadRepository) *GetUserUseCase {
	return &GetUserUseCase{readRepo: readModel}
}

func (h *GetUserUseCase) Execute(ctx context.Context, query GetUserQuery) (*entities.UserRead, error) {
	user, err := h.readRepo.FindByID(ctx, query.TenantID, query.UserID)
	if err != nil {
		return nil, exceptions.ErrUserNotFound
	}

	return user, nil
}
