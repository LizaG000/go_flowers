package router

import (
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

	server.Use(
		gin.Recovery(),
		middleware.RequestLogger(),
	)

	registerFlowerRoutes(server, flowerController)
	registerUserRoutes(server, userController, auth)
	registerLoginRoutes(server, loginController)
	registerFavoriteRoutes(server, favoriteController, auth)

	return server
}
