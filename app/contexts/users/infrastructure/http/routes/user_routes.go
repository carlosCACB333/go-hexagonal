package routes

import (
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/commands"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/queries"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/infrastructure/http/controllers"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/config"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(
	app fiber.Router,
	cfg *config.Config,
	createUseCase *commands.CreateUserUseCase,
	getUseCase *queries.GetUserUseCase,
) {

	users := app.Group("/v1/users")

	users.Post("/",
		middleware.RateLimiterMiddleware(cfg, 10),
		controllers.CreateUserController(createUseCase),
	)

	users.Get("/:id",
		middleware.RateLimiterMiddleware(cfg, 10),
		controllers.GetUserController(getUseCase),
	)
}
