package middleware

import "github.com/gofiber/fiber/v2"

func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Get("X-Tenant-Id")

		if tenantID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "X-Tenant-Id header is required",
			})
		}

		c.Locals("tenant_id", tenantID)
		return c.Next()
	}
}
