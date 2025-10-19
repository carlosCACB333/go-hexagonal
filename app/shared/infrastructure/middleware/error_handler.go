package middleware

import (
	"github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {

		if apiErr, ok := err.(*exceptions.ApiError); ok {
			logger.Warn("Handled error", zap.String("message", apiErr.Message))
			return c.Status(apiErr.Code).JSON(apiErr)
		}

		if fiberErr, ok := err.(*fiber.Error); ok {
			logger.Warn("Fiber error", zap.String("error", fiberErr.Message))
			return c.Status(fiberErr.Code).JSON(exceptions.ApiError{
				Code:    fiberErr.Code,
				Message: fiberErr.Message,
				Detail:  fiberErr.Error(),
			})
		}

		logger.Error("Unhandled error", zap.Error(err))
		internalErr := exceptions.NewInternalServerError("Internal Server Error", err.Error())
		return c.Status(internalErr.Code).JSON(internalErr)
	}
}
