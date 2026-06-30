package rabbitmq

import (
	"fmt"
	"log/slog"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/infra/rabbitmq/handlers"
)

func (rcl *RabbitClient) Consume(
	queueName string,
	logger *slog.Logger,
	flowerHandler *handlers.FlowerHandler,
) error {
	queue, err := rcl.recChan.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("не удалось объявить очередь %q: %w", queueName, err)
	}

	messages, err := rcl.recChan.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("не удалось подписаться на очередь %q: %w", queueName, err)
	}

	logger.Info(
		"ожидание сообщений из RabbitMQ",
		"queue", queue.Name,
	)

	go func() {
		for message := range messages {
			if err := flowerHandler.Create(message); err != nil {
				logger.Error(
					"не удалось обработать RabbitMQ-сообщение",
					"message_id", message.MessageId,
					"error", err,
				)
			}
		}

		logger.Warn(
			"канал получения сообщений RabbitMQ закрыт",
			"queue", queue.Name,
		)
	}()

	return nil
}
