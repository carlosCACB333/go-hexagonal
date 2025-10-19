package events

import (
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	shared_events "github.com/carloscacb333/go-hexagonal/app/shared/domain/events"
)

type UserCreatedEvent struct {
	shared_events.BaseEvent
	Data entities.UserRead `json:"data"`
}

func NewUserCreatedEvent(user *entities.User) UserCreatedEvent {
	return UserCreatedEvent{
		BaseEvent: shared_events.NewBaseEvent("user.created", user.ID.String()),
		Data: entities.UserRead{
			ID:          user.ID,
			TenantID:    user.TenantID,
			Name:        user.Name,
			Email:       user.Email.Value(),
			DisplayName: user.DisplayName,
		},
	}
}
