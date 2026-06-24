package router

import (
	"time"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New(
	auth config.Auth,
	flowerController controller.FlowerController,
	userController controller.UserController,
	loginController controller.LoginController,
	favoriteController controller.FavoriteController,
) *gin.Engine {
	server := gin.New()

	rateLimiter := middleware.NewRateLimiter(
		2,
		time.Minute,
	)

	server.Use(
		gin.Recovery(),
		middleware.RequestLogger(),
		rateLimiter.RateLimiterMiddleware(),
	)
	api := server.Group("/api")

	registerFlowerRoutes(api, flowerController)
	registerUserRoutes(api, userController, auth)
	registerLoginRoutes(api, loginController)
	registerFavoriteRoutes(api, favoriteController, auth)

	return server
}
