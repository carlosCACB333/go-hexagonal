package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/projections"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/events"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
)

type UserProjectionsEventHandler struct {
	userCreatedHandler *projections.UserCreatedHandler
}

func NewUserProjectionsEventHandler(userCreatedHandler *projections.UserCreatedHandler) *UserProjectionsEventHandler {
	return &UserProjectionsEventHandler{
		userCreatedHandler: userCreatedHandler,
	}
}

func (h *UserProjectionsEventHandler) HandleEvent(ctx context.Context, eventType string, data []byte) error {
	switch eventType {
	case "user.created":

		event := &events.UserCreatedEvent{}

		if err := json.Unmarshal(data, event); err != nil {
			return shared_exceptions.NewInternalServerError("failed to unmarshal user.created event", err.Error())
		}

		err := h.userCreatedHandler.Handle(ctx, event)
		return err

	default:
		return shared_exceptions.NewBadRequestError("unknown event type", fmt.Sprintf("event type %s is not recognized", eventType))
	}
}
