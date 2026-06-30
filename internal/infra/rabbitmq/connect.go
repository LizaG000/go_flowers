package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (rcl *RabbitClient) connect(isRec, reconnect bool) (*amqp.Connection, error) {
	if reconnect {
		if isRec {
			rcl.recConn = nil
		} else {
			rcl.sendConn = nil
		}
	}
	c := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		rcl.config.User,
		rcl.config.Password,
		rcl.config.Host,
		rcl.config.Port,
	)

	log.Printf("RabbitMQ DSN: %s", c)

	if isRec && rcl.recConn != nil {
		return rcl.recConn, nil
	}
	if !isRec && rcl.sendConn != nil {
		return rcl.sendConn, nil
	}

	// 🔥 ВСЕГДА ПРАВИЛЬНЫЙ DSN
	c = fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		rcl.config.User,
		rcl.config.Password,
		rcl.config.Host,
		rcl.config.Port,
	)

	log.Printf("RabbitMQ DSN: %s", c)

	conn, err := amqp.Dial(c)
	if err != nil {
		log.Printf("\r\n--- could not create a connection ---\r\n")
		time.Sleep(1 * time.Second)
		return nil, err
	}

	if isRec {
		rcl.recConn = conn
		return rcl.recConn, nil
	}

	rcl.sendConn = conn
	return rcl.sendConn, nil
}
