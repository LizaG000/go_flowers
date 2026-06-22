package router

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"
	"github.com/gin-gonic/gin"
)

func registerFavoriteRoutes(
	router *gin.Engine,
	favoriteController controller.FavoriteController,
) {
	favorites := router.Group("/favorites")

	favorites.POST("", favoriteController.Create)
	favorites.GET("", favoriteController.GetByUserID)
	favorites.DELETE("", favoriteController.Delete)
}
