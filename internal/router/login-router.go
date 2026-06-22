package router

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerLoginRoutes(
	router *gin.Engine,
	loginController controller.LoginController,
) {
	auth := router.Group("/auth")

	auth.POST("/login", loginController.Login)
	auth.POST("/registration", loginController.Registration)
}
