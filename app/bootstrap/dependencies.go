package bootstrap

import (
	"fmt"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/commands"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/queries"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/ports"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/infrastructure/messaging/consumers"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/infrastructure/persistence"
	shared_ports "github.com/carloscacb333/go-hexagonal/app/shared/domain/ports"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/config"
	shared_persistence "github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/persistence"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/rabbitmq"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/security"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	config *config.Config
	logger *zap.Logger
	db     *gorm.DB
	hasher shared_ports.Hasher

	// Event bus
	eventBus shared_ports.EventBus

	// Repositorios
	idempotencyRepository shared_ports.IdempotencyRepository
	userRepository        ports.UserRepository
	userReadRepository    ports.UserReadRepository

	// Casos de uso
	createUserUseCase     *commands.CreateUserUseCase
	getUserUseCase        *queries.GetUserUseCase
	createUserReadUseCase *commands.CreateUserReadUseCase
	sendEmailUseCase      *commands.SendEmailUseCase

	// Consumidores de eventos
	eventConsumers []shared_ports.EventConsumer
}

func NewContainer(cfg *config.Config, logger *zap.Logger, db *gorm.DB) (*Container, error) {
	container := &Container{
		config: cfg,
		logger: logger,
		db:     db,
		hasher: security.NewBcryptHasher(),
	}

	if err := container.initEventBus(); err != nil {
		return nil, fmt.Errorf("failed to initialize event bus: %w", err)
	}

	container.initRepositories()
	container.initUseCases()

	if err := container.initConsumers(); err != nil {
		return nil, fmt.Errorf("failed to initialize consumers: %w", err)
	}

	return container, nil
}

func (c *Container) initEventBus() error {

	eventBus, err := rabbitmq.NewRabbitMQEventBus(&c.config.RabbitMQ)
	if err != nil {
		return fmt.Errorf("failed to create event bus: %w", err)
	}
	c.eventBus = eventBus

	return nil
}

func (c *Container) initRepositories() {
	c.idempotencyRepository = shared_persistence.NewGormIdempotencyRepository(c.db)
	c.userRepository = persistence.NewGormUserRepository(c.db)
	c.userReadRepository = persistence.NewGormUserReadRepository(c.db)

}

func (c *Container) initUseCases() {
	c.createUserUseCase = commands.NewCreateUserUseCase(
		c.userRepository,
		c.idempotencyRepository,
		c.eventBus,
		c.hasher,
	)
	c.getUserUseCase = queries.NewGetUserUseCase(c.userReadRepository)
	c.createUserReadUseCase = commands.NewCreateUserReadUseCase(c.userReadRepository)
	c.sendEmailUseCase = commands.NewSendEmailUseCase()
}

func (c *Container) initConsumers() error {
	userProjectionsConsumer := consumers.NewRabbitMQUserProjectionsConsumer(
		&c.config.RabbitMQ,
		c.logger,
		c.createUserReadUseCase,
	)

	userNotificationConsumer := consumers.NewRabbitMQUserNotificationConsumer(
		&c.config.RabbitMQ,
		c.logger,
		c.sendEmailUseCase,
	)

	c.eventConsumers = []shared_ports.EventConsumer{
		userProjectionsConsumer,
		userNotificationConsumer,
	}

	return nil
}

func (c *Container) GetCreateUserUseCase() *commands.CreateUserUseCase {
	return c.createUserUseCase
}

func (c *Container) GetGetUserUseCase() *queries.GetUserUseCase {
	return c.getUserUseCase
}

func (c *Container) GetSendEmailUseCase() *commands.SendEmailUseCase {
	return c.sendEmailUseCase
}

func (c *Container) GetEventConsumers() []shared_ports.EventConsumer {
	return c.eventConsumers
}

func (c *Container) GetEventBus() shared_ports.EventBus {
	return c.eventBus
}

func (c *Container) GetConfig() *config.Config {
	return c.config
}

func (c *Container) Close() error {
	var errs []error

	if c.eventBus != nil {
		c.eventBus.Close()
	}

	for _, consumer := range c.eventConsumers {
		if consumer != nil {
			if err := consumer.Stop(); err != nil {
				errs = append(errs, fmt.Errorf("failed to stop event consumer: %w", err))
			}
		}
	}

	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get database connection: %w", err))
		} else {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close database: %w", err))
			}
		}
	}

	if c.logger != nil {
		if err := c.logger.Sync(); err != nil {
			c.logger.Debug("failed to sync logger", zap.Error(err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("container close errors: %v", errs)
	}

	return nil
}
