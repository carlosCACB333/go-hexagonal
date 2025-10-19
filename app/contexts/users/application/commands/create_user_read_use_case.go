package commands

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/ports"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
)

type CreateUserReadUseCase struct {
	userReadRepo ports.UserReadRepository
}

func NewCreateUserReadUseCase(userReadRepo ports.UserReadRepository) *CreateUserReadUseCase {
	return &CreateUserReadUseCase{
		userReadRepo: userReadRepo,
	}
}

func (uc *CreateUserReadUseCase) Execute(ctx context.Context, user *entities.UserRead) (*entities.UserRead, error) {

	if err := uc.userReadRepo.Upsert(ctx, user); err != nil {
		return nil, shared_exceptions.NewInternalServerError("failed to create user read model", err.Error())
	}
	return user, nil
}
