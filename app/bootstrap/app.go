package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/infrastructure/http/routes"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
	shared_ports "github.com/carloscacb333/go-hexagonal/app/shared/domain/ports"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/config"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/middleware"
	shared_persistence "github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/persistence"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	container  *Container
	httpServer *fiber.App
	logger     *zap.Logger
}

func NewApplication() (*App, error) {
	// 1. Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Inicializar logger
	logger, err := initLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 3. Conectar base de datos
	db, err := initDatabase(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 4. Construir contenedor de dependencias
	container, err := NewContainerBuilder().
		WithConfig(cfg).
		WithLogger(logger).
		WithDatabase(db).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build container: %w", err)
	}

	// 5. Crear servidor HTTP
	httpServer := createHTTPServer(container, logger)

	return &App{
		container:  container,
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

func initLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func initDatabase(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	db, err := shared_persistence.ConnectDatabase(&cfg.DB)
	if err != nil {
		return nil, err
	}

	// Ejecutar migraciones en desarrollo
	if cfg.App.Environment == "development" {
		logger.Info("running database migrations")
		if err := shared_persistence.AutoMigrate(db, logger); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return db, nil
}

func createHTTPServer(container *Container, logger *zap.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          middleware.ErrorHandler(logger),
		DisableStartupMessage: false,
		AppName:               "Go Hexagonal API",
		ServerHeader:          "Fiber",
		StrictRouting:         true,
		CaseSensitive:         true,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           120 * time.Second,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Request-ID, X-Tenant-ID",
	}))
	app.Use(middleware.CorrelationIDMiddleware())
	app.Use(middleware.TenantMiddleware())
	app.Use(middleware.LoggerMiddleware(logger))

	registerHealthChecks(app, container)

	registerRoutes(app, container)

	return app
}

func registerHealthChecks(app *fiber.App, container *Container) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "go-hexagonal",
		})
	})

	app.Get("/ready", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		// Verificar base de datos
		sqlDB, err := container.db.DB()
		if err != nil {
			e := shared_exceptions.NewServiceUnavailableError(
				"database_unavailable",
				"failed to get database connection",
			)
			return c.Status(e.Code).JSON(e)
		}

		if err := sqlDB.PingContext(ctx); err != nil {
			e := shared_exceptions.NewServiceUnavailableError(
				"database_unreachable",
				"database ping failed",
			)
			return c.Status(e.Code).JSON(e)
		}

		return c.JSON(fiber.Map{
			"status":   "ready",
			"database": "connected",
		})
	})
}

func registerRoutes(app *fiber.App, container *Container) {
	api := app.Group("/api")

	routes.RegisterUserRoutes(
		api,
		container.GetConfig(),
		container.GetCreateUserUseCase(),
		container.GetGetUserUseCase(),
	)
}

func (a *App) StartHTTPServer() error {
	cfg := a.container.GetConfig()
	addr := fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port)

	a.logger.Info("starting HTTP server",
		zap.String("address", addr),
		zap.String("environment", cfg.App.Environment),
	)

	if err := a.httpServer.Listen(addr); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (a *App) StartEventConsumers(ctx context.Context) error {
	a.logger.Info("starting event consumers")
	consumers := a.container.GetEventConsumers()
	errChan := make(chan error, len(consumers))

	for _, consumer := range consumers {
		go func(c shared_ports.EventConsumer) {
			if err := c.Start(ctx); err != nil {
				errChan <- fmt.Errorf("failed to start consumer: %w", err)
			}
		}(consumer)
	}

	for range consumers {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down application")

	if err := a.httpServer.ShutdownWithContext(ctx); err != nil {
		a.logger.Error("error shutting down HTTP server", zap.Error(err))
	}

	if err := a.container.Close(); err != nil {
		a.logger.Error("error closing container", zap.Error(err))
		return err
	}

	a.logger.Info("application shutdown complete")
	return nil
}

func (a *App) GetLogger() *zap.Logger {
	return a.logger
}
