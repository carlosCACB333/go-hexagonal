package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
	"github.com/carloscacb333/go-hexagonal/app/shared/domain/ports"
	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQEventBus struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQEventBus(cfg *config.RabbitMQConfig) (*RabbitMQEventBus, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, exceptions.NewServiceUnavailableError("failed to connect to RabbitMQ", err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, exceptions.NewServiceUnavailableError("failed to open channel", err.Error())
	}

	// Declarar exchange
	if err := channel.ExchangeDeclare(
		"domain_events", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	); err != nil {
		channel.Close()
		conn.Close()
		return nil, exceptions.NewServiceUnavailableError("failed to declare exchange", err.Error())
	}

	return &RabbitMQEventBus{
		conn:    conn,
		channel: channel,
	}, nil
}

func (b *RabbitMQEventBus) Publish(ctx context.Context, event ports.DomainEvent, correlationID string) error {
	data, err := json.Marshal(event)
	if err != nil {
		return exceptions.NewBadRequestError("failed to marshal event", err.Error())
	}

	return b.channel.PublishWithContext(
		ctx,
		"domain_events",   // exchange
		event.EventType(), // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          data,
			DeliveryMode:  amqp.Persistent,
			CorrelationId: correlationID,
			MessageId:     event.EventID(),
		},
	)
}

func (b *RabbitMQEventBus) Close() {
	if b.channel != nil {
		b.channel.Close()
	}
	if b.conn != nil {
		b.conn.Close()
	}
}
