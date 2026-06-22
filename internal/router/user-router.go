package router

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerUserRoutes(
	router *gin.Engine,
	userController controller.UserController,
) {
	user := router.Group("/users")

	user.GET("", userController.Get)
}
