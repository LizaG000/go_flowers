package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"gilab.com/pragmaticrewies/golang-gin-poc/client/dto"
	"gilab.com/pragmaticrewies/golang-gin-poc/client/entity"

	"github.com/google/uuid"
)

func (rcl *RabbitClient) ConsumeResponse(
	queueName string,
	requestID uuid.UUID,
	logger *slog.Logger,
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
		return fmt.Errorf("не удалось объявить очередь ответов %q: %w", queueName, err)
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
		return fmt.Errorf("не удалось подписаться на очередь ответов %q: %w", queueName, err)
	}

	logger.Info(
		"ожидание ответа от RabbitMQ",
		"queue", queue.Name,
		"correlation_id", requestID,
	)

	for message := range messages {
		var response dto.RabbitResponse

		if err := json.Unmarshal(message.Body, &response); err != nil {
			logger.Error(
				"не удалось прочитать RabbitMQ-ответ",
				"message_id", message.MessageId,
				"error", err,
			)

			_ = message.Nack(false, false)
			continue
		}

		if response.CorrelationID != requestID {
			logger.Warn(
				"получен ответ для другого запроса",
				"expected_correlation_id", requestID,
				"received_correlation_id", response.CorrelationID,
			)

			_ = message.Nack(false, true)
			continue
		}

		if response.Status == "error" {
			logger.Error(
				"сервер вернул ошибку",
				"correlation_id", response.CorrelationID,
				"error", response.Error,
			)

			if err := message.Ack(false); err != nil {
				return fmt.Errorf("не удалось подтвердить сообщение с ошибкой: %w", err)
			}

			return nil
		}

		var flower entity.Flower

		if err := json.Unmarshal(response.Data, &flower); err != nil {
			_ = message.Nack(false, false)

			return fmt.Errorf("не удалось прочитать цветок из ответа: %w", err)
		}

		logger.Info(
			"получен ответ о создании цветка",
			"correlation_id", response.CorrelationID,
			"flower_id", flower.ID,
			"title", flower.Title,
			"description", flower.Description,
			"price", flower.Price,
			"height", flower.Height,
			"count", flower.Count,
		)

		if err := message.Ack(false); err != nil {
			return fmt.Errorf("не удалось подтвердить ответ RabbitMQ: %w", err)
		}

		return nil
	}

	return fmt.Errorf("канал получения ответов RabbitMQ закрыт")
}
