package commands

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/events"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/ports"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/value_objects"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
	shared_ports "github.com/carloscacb333/go-hexagonal/app/shared/domain/ports"
	"github.com/google/uuid"
)

type CreateUserCommand struct {
	TenantID       string
	IdempotencyKey string
	CorrelationID  string
	Name           string
	Email          string
	Password       string
	DisplayName    *string
}

type CreateUserResponse struct {
	UserID uuid.UUID `json:"user_id"`
}

type CreateUserUseCase struct {
	userRepo        ports.UserRepository
	idempotencyRepo shared_ports.IdempotencyRepository
	eventBus        shared_ports.EventBus
	hasher          shared_ports.Hasher
}

func NewCreateUserUseCase(
	userRepo ports.UserRepository,
	idempotencyRepo shared_ports.IdempotencyRepository,
	eventBus shared_ports.EventBus,
	hasher shared_ports.Hasher,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:        userRepo,
		idempotencyRepo: idempotencyRepo,
		eventBus:        eventBus,
		hasher:          hasher,
	}
}

func (h *CreateUserUseCase) Execute(ctx context.Context, cmd CreateUserCommand) (*CreateUserResponse, error) {
	// Verificar idempotencia
	if cmd.IdempotencyKey != "" {
		processed, err := h.idempotencyRepo.IsProcessed(ctx, cmd.TenantID, cmd.IdempotencyKey)
		if err != nil {
			return nil, err
		}
		if processed {
			return nil, shared_exceptions.NewConflictError("command already processed", "")
		}
	}

	// Crear value objects
	email, err := value_objects.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	// Verificar email único
	exists, err := h.userRepo.ExistsByEmail(ctx, cmd.TenantID, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, exceptions.ErrDuplicateEmail
	}

	password, err := value_objects.NewPassword(h.hasher, cmd.Password)

	if err != nil {
		return nil, err
	}

	// Crear entidad
	user, err := entities.NewUser(cmd.TenantID, cmd.Name, email, password, cmd.DisplayName)
	if err != nil {
		return nil, err
	}

	// Persistir
	if err := h.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	// Marcar como procesado
	if cmd.IdempotencyKey != "" {
		if err := h.idempotencyRepo.MarkAsProcessed(ctx, cmd.TenantID, cmd.IdempotencyKey); err != nil {
			return nil, err
		}
	}

	// Publicar eventos
	event := events.NewUserCreatedEvent(user)

	if err := h.eventBus.Publish(ctx, event, cmd.CorrelationID); err != nil {
		return nil, err
	}

	return &CreateUserResponse{UserID: user.ID}, nil
}
