package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/dto"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ResponsePublisher interface {
	Publish(queueName string, response dto.RabbitResponse) error
}

type FlowerHandler struct {
	flowerService      service.FlowerService
	idempotencyService service.IdempotencyService
	publisher          ResponsePublisher
	responseQueue      string
	log                *slog.Logger
}

func NewFlowerHandler(
	flowerService service.FlowerService,
	idempotencyService service.IdempotencyService,
	publisher ResponsePublisher,
	responseQueue string,
	log *slog.Logger,
) *FlowerHandler {
	return &FlowerHandler{
		flowerService:      flowerService,
		idempotencyService: idempotencyService,
		publisher:          publisher,
		responseQueue:      responseQueue,
		log:                log,
	}
}

func (h *FlowerHandler) Create(message amqp.Delivery) error {
	var request dto.RabbitRequest

	if err := json.Unmarshal(message.Body, &request); err != nil {
		_ = message.Nack(false, false)

		return fmt.Errorf("не удалось прочитать сообщение RabbitMQ: %w", err)
	}

	if request.Action != "create_flower" {
		_ = message.Nack(false, false)

		return fmt.Errorf("неподдерживаемое действие: %s", request.Action)
	}

	var createFlower entity.CreateFlower

	if err := json.Unmarshal(request.Data, &createFlower); err != nil {
		_ = message.Nack(false, false)

		return fmt.Errorf("не удалось прочитать данные цветка: %w", err)
	}

	token := request.Auth
	idempotencyKey := request.IdempotencyKey

	h.log.Info(
		"получен запрос на создание цветка",
		"message_id", request.ID,
		"action", request.Action,
		"idempotency_key", idempotencyKey,
		"token_received", token != "",
		"title", createFlower.Title,
		"price", createFlower.Price,
		"height", createFlower.Height,
		"count", createFlower.Count,
	)

	createdFlower, err := h.flowerService.Create(createFlower)
	if err != nil {
		_ = message.Nack(false, true)

		return fmt.Errorf("не удалось создать цветок: %w", err)
	}

	flowerBytes, err := json.Marshal(createdFlower)
	if err != nil {
		_ = message.Nack(false, true)

		return fmt.Errorf("не удалось сериализовать созданный цветок: %w", err)
	}

	response := dto.RabbitResponse{
		CorrelationID: request.ID,
		Status:        "ok",
		Data:          flowerBytes,
		Error:         "",
	}

	if err := h.publisher.Publish(h.responseQueue, response); err != nil {
		_ = message.Nack(false, true)

		return fmt.Errorf("не удалось отправить ответ RabbitMQ: %w", err)
	}

	if err := message.Ack(false); err != nil {
		return fmt.Errorf("не удалось подтвердить сообщение RabbitMQ: %w", err)
	}

	h.log.Info(
		"цветок успешно создан, ответ отправлен",
		"message_id", request.ID,
		"flower_id", createdFlower.ID,
		"idempotency_key", idempotencyKey,
		"response_queue", h.responseQueue,
	)

	return nil
}
