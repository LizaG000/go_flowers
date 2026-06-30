package main

import (
	"log/slog"
	"os"

	_ "gilab.com/pragmaticrewies/golang-gin-poc/docs"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/infra/rabbitmq"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/infra/rabbitmq/handlers"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/infra/storage/postgres"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/middleware"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/repository"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/router"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/service"
)

// @title Flower Shop API
// @version 1.0
// @description API интернет-магазина цветов.
// @host localhost:8080
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.MustLoad()

	log := middleware.SetupLogger(cfg.Env)
	slog.SetDefault(log)

	storage, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer storage.DB.Close()

	log.Info("connected to PostgreSQL")
	log.Info("starting flower shop backend", slog.String("env", cfg.Env))

	rabbit, err := rabbitmq.RabbitMQNew(cfg.RabbitMQ, log)
	if err != nil {
		log.Error(
			"failed to connect to RabbitMQ",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer rabbit.Close()

	flowerRepository := repository.NewFlowerRepository(storage.DB)
	userRepository := repository.NewUserRepository(storage.DB)
	passwordRepository := repository.NewPasswordRepository(storage.DB)
	favoriteRepository := repository.NewFavoriteREpository(storage.DB)
	idempotencyRepository := repository.NewIdempotencyRepository(storage.DB)

	flowerService := service.NewFlowerService(flowerRepository)
	userService := service.NewUserService(userRepository)
	loginService := service.NewLoginService(
		storage.DB,
		userRepository,
		passwordRepository,
		cfg.Auth,
	)
	favoriteService := service.NewFavoriteService(favoriteRepository)
	idempotencyService := service.NewIdempotencyService(idempotencyRepository)

	flowerController := controller.NewFlowerController(flowerService)
	userController := controller.NewUserController(userService)
	loginController := controller.NewLoginController(loginService)
	favoriteController := controller.NewFavoriteController(
		favoriteService,
		idempotencyService,
	)
	healthController := controller.NewHealthController(storage.DB)

	flowerHandler := handlers.NewFlowerHandler(
		flowerService,
		idempotencyService,
		rabbit,
		cfg.RabbitMQ.ResponseQueue,
		log,
	)

	if err := rabbit.Consume(
		cfg.RabbitMQ.RequestQueue,
		log,
		flowerHandler,
	); err != nil {
		log.Error(
			"failed to start RabbitMQ consumer",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	server := router.New(
		cfg.Auth,
		flowerController,
		userController,
		loginController,
		favoriteController,
		healthController,
	)

	if err := server.Run(cfg.HTTPServer.Address); err != nil {
		log.Error(
			"failed to start HTTP server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
