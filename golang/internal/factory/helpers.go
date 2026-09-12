package factory

import (
	"fmt"

	"github.com/7574-sistemas-distribuidos/tp-mom/golang/internal/middleware"
	m "github.com/7574-sistemas-distribuidos/tp-mom/golang/internal/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
)

func validateConnectionSetting(c *m.ConnSettings) error {
	if c.Hostname == "" {
		return formatError(EMPTY_CONN_HOST_ERROR)
	}
	if c.Port == 0 {
		return formatError(EMPTY_CONN_HOST_ERROR)
	}
	return nil
}

func buildDialURL(host string, port int) string {
	return fmt.Sprintf("amqp://guest:guest@%s:%d/", host, port)
}

func buildConsummerTag(queueName string) string {
	return fmt.Sprintf("%s-consumer-tag", queueName)
}

func rabbitConnect(c *m.ConnSettings) (*amqp.Connection, error) {
	if err := validateConnectionSetting(c); err != nil {
		return nil, err
	}

	url := buildDialURL(c.Hostname, c.Port)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, formatError(CONN_ERROR)
	}
	return conn, nil
}

func createDefaultRabbitChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, formatError(CHANNEL_CREATE_ERROR)
	}
	return ch, nil
}

func createDefaultRabbitQueue(ch *amqp.Channel, queueName string) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durability
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)

	if err != nil {
		return nil, formatError(QUEUE_CREATE_ERROR)
	}

	return &q, nil
}

// Si ocurre un error interno que no puede resolverse devuelve ErrMessageMiddlewareClose.
func checkRabbitCloseError(err error) error {
	switch err {
	case amqp.ErrClosed:
		return nil
	default:
		return middleware.ErrMessageMiddlewareClose
	}
}

// Si se pierde la conexión con el middleware devuelve ErrMessageMiddlewareDisconnected.
// Si ocurre un error interno que no puede resolverse devuelve ErrMessageMiddlewareMessage.
func checkRabbitSendError(err error) error {
	switch err {
	case amqp.ErrClosed:
		return middleware.ErrMessageMiddlewareDisconnected
	default:
		return middleware.ErrMessageMiddlewareMessage
	}
}

// Si se pierde la conexión con el middleware devuelve ErrMessageMiddlewareDisconnected.
// Si ocurre un error interno que no puede resolverse devuelve ErrMessageMiddlewareMessage.
func checkRabbitConsumError(err error) error {
	switch err {
	case amqp.ErrClosed:
		return middleware.ErrMessageMiddlewareDisconnected
	default:
		return middleware.ErrMessageMiddlewareMessage
	}
}
