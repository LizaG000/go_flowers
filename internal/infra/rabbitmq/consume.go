package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/dto"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/infra/rabbitmq/handlers"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/infra/storage"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/security"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/service"

	"github.com/google/uuid"
)

func (rcl *RabbitClient) Consume(
	queueName string,
	logger *slog.Logger,
	flowerHandler *handlers.FlowerHandler,
	idempotencyService service.IdempotencyService,
	auth config.Auth,
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
		return fmt.Errorf(
			"не удалось объявить очередь %q: %w",
			queueName,
			err,
		)
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
		return fmt.Errorf(
			"не удалось подписаться на очередь %q: %w",
			queueName,
			err,
		)
	}

	logger.Info(
		"ожидание сообщений из RabbitMQ",
		"queue", queue.Name,
	)

	go func() {
		for message := range messages {
			var request dto.RabbitRequest

			if err := json.Unmarshal(message.Body, &request); err != nil {
				logger.Error(
					"не удалось прочитать RabbitRequest",
					"message_id", message.MessageId,
					"error", err,
				)

				_ = message.Nack(false, false)
				continue
			}

			if request.Auth == "" {
				logger.Error(
					"токен авторизации не передан",
					"message_id", request.ID,
				)

				_ = message.Nack(false, false)
				continue
			}

			const bearerPrefix = "Bearer "

			if !strings.HasPrefix(request.Auth, bearerPrefix) {
				logger.Error(
					"неверный формат токена",
					"message_id", request.ID,
				)

				_ = message.Nack(false, false)
				continue
			}

			tokenString := strings.TrimPrefix(
				request.Auth,
				bearerPrefix,
			)

			if tokenString == "" {
				logger.Error(
					"токен авторизации не передан",
					"message_id", request.ID,
				)

				_ = message.Nack(false, false)
				continue
			}

			claims, err := security.ParseToken(
				tokenString,
				auth.PublicKeyPath,
			)

			if err != nil {
				logger.Error(
					"некорректный или просроченный токен",
					"message_id", request.ID,
					"error", err,
				)

				_ = message.Nack(false, false)
				continue
			}

			logger.Info(
				"RabbitMQ-запрос успешно аутентифицирован",
				"message_id", request.ID,
				"user_id", claims.UserID,
			)

			if request.IdempotencyKey == uuid.Nil {
				logger.Error(
					"не передан ключ идемпотентности",
					"message_id", request.ID,
				)

				_ = message.Nack(false, false)
				continue
			}

			payloadHash := security.CalculatePayloadHash(
				request.Data,
			)

			idempotency, err := idempotencyService.Get(
				request.IdempotencyKey,
			)

			if err == nil {

				if idempotency.PayloadHash != payloadHash {
					logger.Error(
						"ключ идемпотентности уже использован для другого запроса",
						"idempotency_key", request.IdempotencyKey,
					)

					_ = message.Nack(false, false)
					continue
				}

				logger.Info(
					"повторный RabbitMQ-запрос, возвращаем ранее созданный цветок",
					"idempotency_key", request.IdempotencyKey,
					"message_id", request.ID,
				)

				response := dto.RabbitResponse{
					CorrelationID: request.ID,
					Status:        "ok",
					Data:          idempotency.ResponseBody,
					Error:         "",
				}

				if err := flowerHandler.PublishResponse(response); err != nil {
					logger.Error(
						"не удалось отправить сохранённый ответ RabbitMQ",
						"idempotency_key", request.IdempotencyKey,
						"error", err,
					)

					_ = message.Nack(false, true)
					continue
				}

				if err := message.Ack(false); err != nil {
					logger.Error(
						"не удалось подтвердить повторное RabbitMQ-сообщение",
						"idempotency_key", request.IdempotencyKey,
						"error", err,
					)

					continue
				}

				logger.Info(
					"ранее созданный цветок повторно отправлен клиенту",
					"idempotency_key", request.IdempotencyKey,
					"correlation_id", request.ID,
				)

				continue
			}

			if !errors.Is(err, storage.ErrIdempotencyNotFound) {
				logger.Error(
					"не удалось проверить ключ идемпотентности",
					"idempotency_key", request.IdempotencyKey,
					"error", err,
				)

				_ = message.Nack(false, true)
				continue
			}

			flower, err := flowerHandler.Create(message)
			if err != nil {
				logger.Error(
					"не удалось обработать RabbitMQ-сообщение",
					"message_id", message.MessageId,
					"error", err,
				)

				continue
			}

			responseBody, err := json.Marshal(flower)
			if err != nil {
				logger.Error(
					"не удалось сериализовать созданный цветок",
					"idempotency_key", request.IdempotencyKey,
					"error", err,
				)

				continue
			}

			_, err = idempotencyService.Create(
				entity.CreateIdempotency{
					Key:          request.IdempotencyKey,
					Status:       "completed",
					ResponseCode: http.StatusCreated,
					ResponseBody: responseBody,
					PayloadHash:  payloadHash,
				},
			)

			if err != nil {
				logger.Error(
					"не удалось сохранить ключ идемпотентности",
					"idempotency_key", request.IdempotencyKey,
					"error", err,
				)

				continue
			}

			logger.Info(
				"ключ идемпотентности сохранён",
				"idempotency_key", request.IdempotencyKey,
				"flower_id", flower.ID,
			)
		}

		logger.Warn(
			"канал получения сообщений RabbitMQ закрыт",
			"queue", queue.Name,
		)
	}()

	return nil
}
