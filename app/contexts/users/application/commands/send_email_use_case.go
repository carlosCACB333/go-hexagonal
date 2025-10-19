package commands

import (
	"context"
	"fmt"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
)

type SendEmailUseCase struct {
}

func NewSendEmailUseCase() *SendEmailUseCase {
	return &SendEmailUseCase{}
}

func (uc *SendEmailUseCase) Execute(ctx context.Context, user *entities.UserRead) error {

	fmt.Printf("Sending welcome email to %s at %s\n", user.Name, user.Email)
	return nil
}
