package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/commands"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/events"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
)

type UserProjectionsEventHandler struct {
	createUserUseCase *commands.CreateUserReadUseCase
}

func NewUserProjectionsEventHandler(createUserUseCase *commands.CreateUserReadUseCase) *UserProjectionsEventHandler {
	return &UserProjectionsEventHandler{
		createUserUseCase: createUserUseCase,
	}
}

func (h *UserProjectionsEventHandler) HandleEvent(ctx context.Context, eventType string, data []byte) error {
	switch eventType {
	case "user.created":

		event := events.UserCreatedEvent{}

		if err := json.Unmarshal(data, &event); err != nil {
			return shared_exceptions.NewInternalServerError("failed to unmarshal user.created event", err.Error())
		}

		_, err := h.createUserUseCase.Execute(ctx, &event.Data)
		return err

	default:
		return shared_exceptions.NewBadRequestError("unknown event type", fmt.Sprintf("event type %s is not recognized", eventType))
	}
}
