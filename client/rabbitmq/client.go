package rabbitmq

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/client/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	sendConn *amqp.Connection
	recConn  *amqp.Connection
	sendChan *amqp.Channel
	recChan  *amqp.Channel

	config config.RabbitMQ
}
