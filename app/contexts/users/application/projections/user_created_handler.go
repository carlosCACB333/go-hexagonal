package projections

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/events"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/ports"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
)

type UserCreatedHandler struct {
	userReadRepo ports.UserReadRepository
}

func NewUserCreatedHandler(userReadRepo ports.UserReadRepository) *UserCreatedHandler {
	return &UserCreatedHandler{
		userReadRepo: userReadRepo,
	}
}

func (uc *UserCreatedHandler) Handle(ctx context.Context, event *events.UserCreatedEvent) error {
	if event == nil {
		return shared_exceptions.NewInternalServerError("failed to create user read model", "event cannot be nil")
	}

	if err := uc.userReadRepo.Upsert(ctx, &event.Data); err != nil {
		return shared_exceptions.NewInternalServerError("failed to create user read model", err.Error())
	}
	return nil
}
