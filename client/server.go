package main

import (
	"log"
	"log/slog"
	"os"

	"gilab.com/pragmaticrewies/golang-gin-poc/client/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/client/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/client/rabbitmq"
	"github.com/google/uuid"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, nil),
	)
	var idempotencyKey uuid.UUID

	idempotencyKey, err := uuid.Parse("7f3c2b91-6d44-4a8e-9f12-b5c7d0e36a21")
	if err != nil {
		log.Fatal(err)
	}

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

	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjlkMzBjNWMtNzA5Ny00MGM5LWE5YjYtNDc0MGYzY2ZlMDU0IiwiZXhwIjoxNzkwNzc5MzM4LCJpYXQiOjE3OTA2OTI5Mzh9.iY3eu69TbCu6ZvBpEDANGpY2kr86uSqjK1Q57WL_smsdPr2oH3IG29A0jI81x2pcrAmg4gIHma1GbmmXpCJX0SXs7bSATEEuBHJf9eFARXiB53uyQjlg6GKScbJxjfh_QicVacKAU6_pJUwQEnLqvcGBrHMZUZHq4xXNmJr-eYO0qr_ybOTuGRYTU1RnIn9QhTOmdDs41MtpGoWw9A2SrenJDZiwpLkLsVufH6DrCurphrH4lMMO9utIKTfnNq67dvu6HuHZP-ZiNweGsdoKasxxRD3wmoDCgiXqaO6pHaplMZ4qQcesdN2RDWutSNvl-W-GlQdO91uHhIf45d1ZVA"

	requestID, err := rabbitClient.Publish(
		cfg.RabbitMQ.RequestQueue,
		flower,
		"Bearer "+jwtToken,
		idempotencyKey,
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
