package controllers

import (
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/commands"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
	"github.com/gofiber/fiber/v2"
)

type CreateUserRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8"`
	DisplayName *string `json:"display_name,omitempty"`
}

func CreateUserController(useCase *commands.CreateUserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateUserRequest

		if err := c.BodyParser(&req); err != nil {
			return shared_exceptions.NewBadRequestError("invalid request body", err.Error())
		}

		tenantID := c.Locals("tenant_id").(string)
		correlationID := c.Locals("correlation_id").(string)
		idempotencyKey := c.Get("X-Idempotency-Key")

		// Feature flag: display_name
		if req.DisplayName != nil && !isFeatureEnabled(c, "display_name") {
			req.DisplayName = nil
		}

		cmd := commands.CreateUserCommand{
			TenantID:       tenantID,
			IdempotencyKey: idempotencyKey,
			CorrelationID:  correlationID,
			Name:           req.Name,
			Email:          req.Email,
			Password:       req.Password,
			DisplayName:    req.DisplayName,
		}

		resp, err := useCase.Execute(c.Context(), cmd)
		if err != nil {
			return err
		}

		return c.Status(201).JSON(fiber.Map{
			"user_id": resp.UserID,
			"message": "User created successfully",
		})
	}
}

func isFeatureEnabled(c *fiber.Ctx, feature string) bool {
	// Implementar lógica de feature flags por tenant
	// Por ahora retorna true
	return true
}
