package router

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"

	"github.com/gin-gonic/gin"
)

func registerFlowerRoutes(
	router *gin.Engine,
	flowerController controller.FlowerController,
) {
	flowers := router.Group("/flowers")

	flowers.GET("", flowerController.GetAll)
	flowers.POST("", flowerController.Create)
	flowers.PUT("", flowerController.Update)
	flowers.DELETE("", flowerController.Delete)
}
