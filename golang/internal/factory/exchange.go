package factory

import (
	"github.com/7574-sistemas-distribuidos/tp-mom/golang/internal/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Exchange struct {
	name         string
	connection   *amqp.Connection
	channel      *amqp.Channel
	consummerTag *string
}

func (q *Exchange) StartConsuming(callbackFunc func(msg middleware.Message, ack func(), nack func())) error {
	return nil
}

func (q *Exchange) StopConsuming() error {
	return nil
}

func (q *Exchange) Send(msg middleware.Message) error {
	return nil
}

func (q *Exchange) Close() error {
	if q.channel != nil {
		if err := q.channel.Close(); err != nil {
			return checkRabbitCloseError(err)
		}
	}
	if q.connection != nil {
		if err := q.connection.Close(); err != nil {
			return checkRabbitCloseError(err)
		}
	}
	return nil
}
