package main

import (
	"log"
	"log/slog"
	"os"

	"gilab.com/pragmaticrewies/golang-gin-poc/client/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/client/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/client/rabbitmq"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, nil),
	)

	cfg := config.MustLoad()

	rabbitClient, err := rabbitmq.RabbitMQNew(
		cfg.RabbitMQ,
		logger,
	)
	if err != nil {
		log.Fatal("не удалось подключиться к RabbitMQ: ", err)
	}
	defer rabbitClient.Close()

	flower := entity.CreateFlower{
		Title:       "Красная роза",
		Description: "Свежая красная роза",
		Price:       250,
		Height:      50,
		Count:       10,
	}

	jwtToken := ""

	requestID, err := rabbitClient.Publish(
		cfg.RabbitMQ.RequestQueue,
		flower,
		"Bearer "+jwtToken,
	)
	if err != nil {
		log.Fatal("не удалось отправить сообщение: ", err)
	}

	logger.Info(
		"запрос на создание цветка отправлен",
		"queue", cfg.RabbitMQ.RequestQueue,
		"request_id", requestID,
	)

	if err := rabbitClient.ConsumeResponse(
		cfg.RabbitMQ.ResponseQueue,
		requestID,
		logger,
	); err != nil {
		log.Fatal("не удалось получить ответ RabbitMQ: ", err)
	}
}
