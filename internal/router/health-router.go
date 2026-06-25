package router

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerHealthRoutes(
	router *gin.RouterGroup,
	healthController controller.HealthController,
) {
	flowers := router.Group("/health")

	flowers.GET("", healthController.Check)
}
