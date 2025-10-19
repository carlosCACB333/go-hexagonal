package consumers

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/commands"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/infrastructure/messaging/handlers"
	"github.com/carloscacb333/go-hexagonal/app/shared/domain/ports"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/config"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/rabbitmq"
	"go.uber.org/zap"
)

type RabbitMQUserProjectionsConsumer struct {
	consumer ports.EventConsumer
}

func NewRabbitMQUserProjectionsConsumer(
	cfg *config.RabbitMQConfig,
	logger *zap.Logger,
	createUseCase *commands.CreateUserReadUseCase,
) *RabbitMQUserProjectionsConsumer {

	eventHandler := handlers.NewUserProjectionsEventHandler(createUseCase)

	consumer := rabbitmq.NewRabbitMQConsumer(
		cfg,
		logger,
		"domain_events",
		"user_projections",
		[]string{"user.created"},
		eventHandler,
	)

	return &RabbitMQUserProjectionsConsumer{
		consumer: consumer,
	}
}

func (c *RabbitMQUserProjectionsConsumer) Start(ctx context.Context) error {
	return c.consumer.Start(ctx)
}

func (c *RabbitMQUserProjectionsConsumer) Stop() error {
	return c.consumer.Stop()
}
