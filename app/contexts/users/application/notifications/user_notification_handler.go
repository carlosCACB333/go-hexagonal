package notifications

import (
	"context"
	"fmt"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/events"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/ports"
)

type UserNotificationHandler struct {
	userReadRepo ports.UserReadRepository
}

func NewUserNotificationHandler(userReadRepo ports.UserReadRepository) *UserNotificationHandler {
	return &UserNotificationHandler{
		userReadRepo: userReadRepo,
	}
}

func (uc *UserNotificationHandler) Handle(ctx context.Context, event *events.UserCreatedEvent) error {

	fmt.Printf("Sending welcome email to %s at %s\n", event.Data.Name, event.Data.Email)
	return nil
}
