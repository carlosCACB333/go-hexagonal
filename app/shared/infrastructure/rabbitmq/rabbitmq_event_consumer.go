package rabbitmq

import (
	"context"
	"fmt"

	"github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
	"github.com/carloscacb333/go-hexagonal/app/shared/domain/ports"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQConsumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	cfg        *config.RabbitMQConfig
	exchange   string
	queueName  string
	eventTypes []string
	handler    ports.EventHandler
	logger     *zap.Logger
}

func NewRabbitMQConsumer(
	cfg *config.RabbitMQConfig,
	logger *zap.Logger,
	exchange string,
	queueName string,
	eventTypes []string,
	handler ports.EventHandler,
) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		cfg:        cfg,
		logger:     logger,
		exchange:   exchange,
		queueName:  queueName,
		eventTypes: eventTypes,
		handler:    handler,
	}
}

func (c *RabbitMQConsumer) Start(ctx context.Context) error {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s%s",
		c.cfg.User,
		c.cfg.Password,
		c.cfg.Host,
		c.cfg.Port,
		c.cfg.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return exceptions.NewServiceUnavailableError("failed to connect to RabbitMQ", err.Error())
	}
	c.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		return exceptions.NewServiceUnavailableError("failed to open channel", err.Error())
	}
	c.channel = channel

	queue, err := channel.QueueDeclare(
		c.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return exceptions.NewServiceUnavailableError("failed to declare queue", err.Error())
	}

	for _, eventType := range c.eventTypes {
		if err := channel.QueueBind(
			queue.Name,
			eventType,
			c.exchange,
			false,
			nil,
		); err != nil {
			return exceptions.NewServiceUnavailableError("failed to bind queue", err.Error())
		}
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return exceptions.NewServiceUnavailableError("failed to register consumer", err.Error())
	}

	c.logger.Info("RabbitMQ consumer started", zap.String("queue", c.queueName))

	for {
		select {
		case <-ctx.Done():
			return c.Stop()
		case msg := <-msgs:
			if err := c.processMessage(ctx, msg); err != nil {
				msg.Nack(false, true) // requeue
				c.logger.Error("x Error handling message", zap.Error(err))
			} else {
				msg.Ack(false)
				c.logger.Info("✓ Message processed successfully")
			}
		}
	}
}

func (c *RabbitMQConsumer) processMessage(ctx context.Context, msg amqp.Delivery) error {

	eventType := msg.RoutingKey
	eventQueue := c.queueName
	correlationID := msg.CorrelationId
	messageID := msg.MessageId

	c.logger.Info("📩Received event",
		zap.String("queue", eventQueue),
		zap.String("eventType", eventType),
		zap.String("correlationID", correlationID),
		zap.String("messageID", messageID),
		zap.ByteString("body", msg.Body))

	return c.handler.HandleEvent(ctx, eventType, msg.Body)
}

func (c *RabbitMQConsumer) Stop() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return exceptions.NewInternalServerError("error closing channel", err.Error())
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return exceptions.NewInternalServerError("error closing connection", err.Error())
		}
	}

	return nil
}
