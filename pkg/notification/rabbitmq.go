package notification

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	User         string `env:"RABBITMQ_USER" envDefault:"guest"`
	Password     string `env:"RABBITMQ_PASSWORD" envDefault:"guest"`
	Host         string `env:"RABBITMQ_HOST" envDefault:"localhost"`
	Port         string `env:"RABBITMQ_PORT" envDefault:"5672"`
	QueueName    string `env:"RABBITMQ_NOTIFICATION_QUEUE_NAME" envDefault:"push-notification-queue"`
	ExchangeName string `env:"RABBITMQ_NOTIFICATION_EXCHANGE_NAME" envDefault:"notification"`
	RoutingKey   string `env:"RABBITMQ_NOTIFICATION_ROUTING_KEY" envDefault:"push-notification"`
}

func (c Config) GetURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", c.User, c.Password, c.Host, c.Port)
}

type RabbitMQNotificationClient struct {
	exchangeName string
	routingKey   string
	connection   *amqp.Connection
	channel      *amqp.Channel
	queue        *amqp.Queue
}

func NewRabbitMQNotificationClient(c *Config) (*RabbitMQNotificationClient, error) {
	conn, err := amqp.Dial(c.GetURL())
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := ch.ExchangeDeclare(c.ExchangeName, "direct", true, false, false, false, nil); err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(c.QueueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	if err := ch.QueueBind(q.Name, c.RoutingKey, c.ExchangeName, false, nil); err != nil {
		return nil, err
	}

	return &RabbitMQNotificationClient{channel: ch, queue: &q, connection: conn, exchangeName: c.ExchangeName, routingKey: c.RoutingKey}, nil
}

func (c *RabbitMQNotificationClient) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.connection.Close()
}

type payload struct {
	ID    string            `json:"id,omitempty"`
	Title string            `json:"title,omitempty"`
	Body  string            `json:"body,omitempty"`
	Data  map[string]string `json:"data,omitempty"`
}

func parseSendParamsToPayload(params SendParams, data map[string]string) *payload {
	return &payload{
		ID:    params.ID,
		Title: params.Title,
		Body:  params.Body,
		Data:  data,
	}
}

func (c *RabbitMQNotificationClient) Send(ctx context.Context, params SendParams, data map[string]string) error {
	payload := parseSendParamsToPayload(params, data)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Priority:     0,
		Body:         payloadJSON,
	}
	return c.channel.PublishWithContext(ctx, c.exchangeName, c.routingKey, false, false, msg)
}
