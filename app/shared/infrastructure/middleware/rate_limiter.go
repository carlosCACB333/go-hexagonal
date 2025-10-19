package middleware

import (
	"sync"
	"time"

	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/config"
	"github.com/gofiber/fiber/v2"
)

type RateLimiter struct {
	requests map[string]*tenantLimiter
	mu       sync.RWMutex
}

type tenantLimiter struct {
	count     int
	resetTime time.Time
}

func RateLimiterMiddleware(cfg *config.Config, maxRequests int) fiber.Handler {
	limiter := &RateLimiter{
		requests: make(map[string]*tenantLimiter),
	}

	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenant_id").(string)

		limiter.mu.Lock()

		now := time.Now()
		tl, exists := limiter.requests[tenantID]

		if !exists || now.After(tl.resetTime) {
			limiter.requests[tenantID] = &tenantLimiter{
				count:     1,
				resetTime: now.Add(time.Minute),
			}
			limiter.mu.Unlock()
			return c.Next()
		}

		if tl.count >= maxRequests {
			limiter.mu.Unlock()
			return c.Status(429).JSON(fiber.Map{
				"error": "rate limit exceeded",
			})
		}

		tl.count++
		limiter.mu.Unlock()

		return c.Next()
	}
}
