package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"gilab.com/pragmaticrewies/golang-gin-poc/client/dto"
	"gilab.com/pragmaticrewies/golang-gin-poc/client/entity"
)

func (rcl *RabbitClient) Publish(
	queueName string,
	flower entity.CreateFlower,
	authToken string,
) (uuid.UUID, error) {
	queue, err := rcl.sendChan.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"не удалось объявить очередь %s: %w",
			queueName,
			err,
		)
	}

	flowerBytes, err := json.Marshal(flower)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"не удалось сериализовать цветок: %w",
			err,
		)
	}

	request := dto.RabbitRequest{
		ID:             uuid.New(),
		Version:        "v1",
		Action:         "create_flower",
		Data:           flowerBytes,
		Auth:           authToken,
		IdempotencyKey: uuid.New(),
	}

	body, err := json.Marshal(request)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"не удалось сериализовать RabbitRequest: %w",
			err,
		)
	}

	err = rcl.sendChan.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			MessageId:     request.ID.String(),
			CorrelationId: request.ID.String(),
			DeliveryMode:  amqp.Persistent,
			ContentType:   "application/json",
			Body:          body,
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"не удалось отправить сообщение в очередь %s: %w",
			queueName,
			err,
		)
	}

	slog.Info(
		"сообщение отправлено в RabbitMQ",
		"queue", queueName,
		"request_id", request.ID,
	)

	return request.ID, nil
}
