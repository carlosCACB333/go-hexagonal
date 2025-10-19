package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/commands"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/events"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
)

type UserNotificationEventHandler struct {
	sendEmailUseCase *commands.SendEmailUseCase
}

func NewUserNotificationEventHandler(sendEmail *commands.SendEmailUseCase) *UserNotificationEventHandler {
	return &UserNotificationEventHandler{
		sendEmailUseCase: sendEmail,
	}
}

func (h *UserNotificationEventHandler) HandleEvent(ctx context.Context, eventType string, data []byte) error {
	switch eventType {
	case "user.created":

		event := events.UserCreatedEvent{}

		if err := json.Unmarshal(data, &event); err != nil {
			return shared_exceptions.NewInternalServerError("failed to unmarshal user.created event", err.Error())
		}

		err := h.sendEmailUseCase.Execute(ctx, &event.Data)
		return err

	default:
		return shared_exceptions.NewBadRequestError("unknown event type", fmt.Sprintf("event type %s is not recognized", eventType))
	}
}
