package controllers

import (
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/queries"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUserController(handler *queries.GetUserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return exceptions.ErrInvalidUuid
		}

		tenantID := c.Locals("tenant_id").(string)

		query := queries.GetUserQuery{
			TenantID: tenantID,
			UserID:   userID,
		}

		user, err := handler.Execute(c.Context(), query)
		if err != nil {
			return err
		}

		return c.JSON(user)
	}
}
