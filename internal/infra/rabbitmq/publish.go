package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/dto"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (rcl *RabbitClient) Publish(
	queueName string,
	response dto.RabbitResponse,
) error {
	queue, err := rcl.sendChan.QueueDeclare(
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

	body, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать RabbitResponse: %w", err)
	}

	err = rcl.sendChan.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			MessageId:     response.CorrelationID.String(),
			CorrelationId: response.CorrelationID.String(),
			DeliveryMode:  amqp.Persistent,
			ContentType:   "application/json",
			Body:          body,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"не удалось отправить сообщение в очередь %q: %w",
			queueName,
			err,
		)
	}

	slog.Info(
		"ответ отправлен в RabbitMQ",
		"queue", queueName,
		"correlation_id", response.CorrelationID,
		"status", response.Status,
	)

	return nil
}
