package rabbitmq

import (
	"fmt"
	"log/slog"

	"gilab.com/pragmaticrewies/golang-gin-poc/client/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RabbitMQNew(
	cfg config.RabbitMQ,
	log *slog.Logger,
) (*RabbitClient, error) {

	// build DSN
	rabbitURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	// 1 connection (правильно)
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к RabbitMQ: %w", err)
	}

	// send channel
	sendChan, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("не удалось открыть send channel: %w", err)
	}

	// receive channel
	recChan, err := conn.Channel()
	if err != nil {
		_ = sendChan.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("не удалось открыть rec channel: %w", err)
	}

	log.Info("подключение к RabbitMQ установлено")

	return &RabbitClient{
		sendConn: conn,
		recConn:  conn,
		sendChan: sendChan,
		recChan:  recChan,
		config:   cfg,
	}, nil
}
