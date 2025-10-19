package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CorrelationIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		correlationID := c.Get("X-Correlation-Id")

		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		c.Locals("correlation_id", correlationID)
		c.Set("X-Correlation-Id", correlationID)

		return c.Next()
	}
}
