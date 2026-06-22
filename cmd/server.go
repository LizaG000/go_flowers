package main

import (
	"log/slog"
	"os"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/middlware"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/repository"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/router"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/service"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()

	log := middlware.SetupLogger(cfg.Env)

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

	flowerRepository := repository.NewFlowerRepository(storage.DB)
	userRepository := repository.NewUserRepository(storage.DB)
	passwordRepository := repository.NewPasswordRepository(storage.DB)

	flowerService := service.NewFlowerService(flowerRepository)
	userService := service.NewUserService(userRepository)
	loginService := service.NewLoginService(storage.DB, userRepository, passwordRepository, cfg.Auth)

	flowerController := controller.NewFlowerController(flowerService)
	userController := controller.NewUserController(userService)
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
