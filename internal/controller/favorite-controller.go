package controller

import (
	"net/http"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/dto"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FavoriteController interface {
	Create(ctx *gin.Context)
	GetByUserID(ctx *gin.Context)
	GetByUserIDOld(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type favoriteController struct {
	service service.FavoriteService
}

func NewFavoriteController(service service.FavoriteService) FavoriteController {
	return &favoriteController{
		service: service,
	}
}

func (c *favoriteController) Create(ctx *gin.Context) {
	var createFavorite dto.RequestCreateFavorite

	if err := ctx.ShouldBindJSON(&createFavorite); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось прочитать данные",
		})
		return
	}

	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "пользователь не авторизован",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось получить id пользователя",
		})
		return
	}

	favorite, err := c.service.Create(entity.CreateFavorite{
		UserID:   userID,
		FlowerID: createFavorite.FlowerID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, favorite)
}

func (c *favoriteController) GetByUserID(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "пользователь не авторизован",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось получить id пользователя",
		})
		return
	}

	favorites, err := c.service.GetByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, favorites)
}

func (c *favoriteController) GetByUserIDOld(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "пользователь не авторизован",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось получить id пользователя",
		})
		return
	}

	favorites, err := c.service.GetByUserIDOld(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, favorites)
}

func (c *favoriteController) Delete(ctx *gin.Context) {
	favoriteIDRaw := ctx.Query("favoriteID")

	if favoriteIDRaw == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id пользователя обязателен",
		})
		return
	}
	favoriteID, err := uuid.Parse(favoriteIDRaw)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	favorites, err := c.service.Delete(favoriteID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, favorites)
}
