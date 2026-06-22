package main

import (
	"log/slog"
	"os"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/router"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/service"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	storage, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer storage.DB.Close()

	slog.SetDefault(log)

	log.Info("connected to PostgreSQL")
	log.Info("starting flower shop backend", slog.String("env", cfg.Env))

	log.Debug("debug messages are enabled")

	flowerRepository := postgres.NewFlowerRepository(storage.DB)
	userRepository := postgres.NewUserRepository(storage.DB)
	passwordRepository := postgres.NewPasswordRepository(storage.DB)

	flowerService := service.NewFlowerService(flowerRepository)
	userService := service.NewUserService(
		storage.DB,
		userRepository,
		passwordRepository,
	)
	loginService := service.NewLoginService(userRepository, passwordRepository, cfg.Auth)

	flowerController := controller.NewFlowerController(flowerService)
	userController := controller.NewUserController(userService, loginService)
	loginController := controller.NewLoginController(loginService)

	server := router.New(
		flowerController,
		userController,
		loginController,
	)
	if err := server.Run(cfg.HTTPServer.Address); err != nil {
		log.Error("failed to start HTTP server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:

		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}
