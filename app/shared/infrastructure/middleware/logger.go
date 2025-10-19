package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) fiber.Handler {

	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		logger.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)),
			zap.String("tenant_id", c.Locals("tenant_id").(string)),
			zap.String("correlation_id", c.Locals("correlation_id").(string)),
		)

		return err
	}
}
