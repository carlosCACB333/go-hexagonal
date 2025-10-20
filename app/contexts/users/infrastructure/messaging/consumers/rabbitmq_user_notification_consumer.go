package consumers

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/application/notifications"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/infrastructure/messaging/handlers"
	"github.com/carloscacb333/go-hexagonal/app/shared/domain/ports"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/config"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/rabbitmq"
	"go.uber.org/zap"
)

type RabbitMQUserNotificationConsumer struct {
	consumer ports.EventConsumer
}

func NewRabbitMQUserNotificationConsumer(
	cfg *config.RabbitMQConfig,
	logger *zap.Logger,
	notificationHandler *notifications.UserNotificationHandler,
) *RabbitMQUserNotificationConsumer {

	eventHandler := handlers.NewUserNotificationEventHandler(notificationHandler)

	consumer := rabbitmq.NewRabbitMQConsumer(
		cfg,
		logger,
		"domain_events",
		"user_notifications",
		[]string{"user.created"},
		eventHandler,
	)

	return &RabbitMQUserNotificationConsumer{
		consumer: consumer,
	}
}

func (c *RabbitMQUserNotificationConsumer) Start(ctx context.Context) error {
	return c.consumer.Start(ctx)
}

func (c *RabbitMQUserNotificationConsumer) Stop() error {
	return c.consumer.Stop()
}
