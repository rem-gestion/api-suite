package broker

import (
	"fmt"

	"github.com/rem-gestion/rem-common/config"
	"github.com/streadway/amqp"
)

// NewRabbit abre la conexión a RabbitMQ usando cfg.RabbitConfig
func NewRabbit(cfg config.RabbitConfig) (*amqp.Connection, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("broker: falta REM_RABBIT_HOST")
	}
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("broker: no pude conectar a RabbitMQ en %s: %w", url, err)
	}
	return conn, nil
}

// OpenChannel abre un canal desde la conexión
func OpenChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	return conn.Channel()
}

// DeclareExchange crea un exchange (tipo topic/direct/fanout)
func DeclareExchange(
	ch *amqp.Channel,
	name, kind string,
	durable, autoDelete, internal, noWait bool,
	args amqp.Table,
) error {
	return ch.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

// DeclareQueue abre o crea una cola
func DeclareQueue(
	ch *amqp.Channel,
	name string,
	durable, autoDelete, exclusive, noWait bool,
	args amqp.Table,
) (amqp.Queue, error) {
	return ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

// BindQueue liga una cola a un exchange con routingKey
func BindQueue(
	ch *amqp.Channel,
	queue, routingKey, exchange string,
	noWait bool,
	args amqp.Table,
) error {
	return ch.QueueBind(queue, routingKey, exchange, noWait, args)
}

// Publish manda un mensaje al exchange con routingKey
func Publish(
	ch *amqp.Channel,
	exchange, routingKey string,
	msg amqp.Publishing,
) error {
	return ch.Publish(exchange, routingKey, false, false, msg)
}

// Consume retorna el canal de entregas (deliveries)
func Consume(
	ch *amqp.Channel,
	queue, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table,
) (<-chan amqp.Delivery, error) {
	return ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}
