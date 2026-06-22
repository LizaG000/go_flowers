package router

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerUserRoutes(
	router *gin.RouterGroup,
	userController controller.UserController,
	auth config.Auth,
) {
	user := router.Group("/users")

	user.GET("", middleware.Auth(auth.PublicKeyPath), userController.Get)
}
