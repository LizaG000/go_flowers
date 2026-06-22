package router

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"

	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

func New(
	flowerController controller.FlowerController,
	userController controller.UserController,
	loginController controller.LoginController,
) *gin.Engine {
	server := gin.New()

	server.Use(
		gin.Recovery(),
		gindump.Dump(),
	)

	registerFlowerRoutes(server, flowerController)
	registerUserRoutes(server, userController)
	registerLoginRoutes(server, loginController)

	return server
}
