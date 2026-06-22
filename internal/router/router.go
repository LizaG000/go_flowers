package router

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"

	"github.com/gin-gonic/gin"
)

func New(
	flowerController controller.FlowerController,
	userController controller.UserController,
	loginController controller.LoginController,
	favoriteController controller.FavoriteController,
) *gin.Engine {
	server := gin.New()

	server.Use(
		gin.Recovery(),
	)

	registerFlowerRoutes(server, flowerController)
	registerUserRoutes(server, userController)
	registerLoginRoutes(server, loginController)
	registerFavoriteRoutes(server, favoriteController)

	return server
}
