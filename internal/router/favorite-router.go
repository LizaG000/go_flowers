package router

import (
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/config"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/controller"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/middleware"
	"github.com/gin-gonic/gin"
)

func registerFavoriteRoutes(
	router *gin.RouterGroup,
	favoriteController controller.FavoriteController,
	auth config.Auth,
) {
	favorites_v1 := router.Group("/v1/favorites")
	favorites_v1.Use(middleware.Auth(auth.PublicKeyPath))
	favorites_v1.GET("", favoriteController.GetByUserIDOld)

	favorites_v2 := router.Group("/v2/favorites")
	favorites_v2.Use(middleware.Auth(auth.PublicKeyPath))
	favorites_v2.GET("", favoriteController.GetByUserID)

	favorites := router.Group("/favorites")
	favorites.Use(middleware.Auth(auth.PublicKeyPath))
	favorites.POST("", favoriteController.Create)
	favorites.DELETE("", favoriteController.Delete)
}
