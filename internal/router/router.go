package router

import (
	"time"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		12,
		time.Minute,
	)

	server.Use(
		gin.Recovery(),
		middleware.RequestLogger(),
		rateLimiter.RateLimiterMiddleware(),
	)

	server.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	api := server.Group("/api")

	registerFlowerRoutes(api, flowerController)
	registerUserRoutes(api, userController, auth)
	registerLoginRoutes(api, loginController)
	registerFavoriteRoutes(api, favoriteController, auth)

	return server
}
