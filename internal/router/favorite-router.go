package router

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/middleware"
	"github.com/gin-gonic/gin"
)

func registerFavoriteRoutes(
	router *gin.Engine,
	favoriteController controller.FavoriteController,
	auth config.Auth,
) {
	favorites := router.Group("/favorites")

	favorites.Use(middleware.Auth(auth.PublicKeyPath))

	favorites.POST("", favoriteController.Create)
	favorites.GET("", favoriteController.GetByUserID)
	favorites.DELETE("", favoriteController.Delete)
}
