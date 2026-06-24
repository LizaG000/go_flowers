package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/dto"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/security"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/service"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/storage"
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
	favoriteService    service.FavoriteService
	idempotencyService service.IdempotencyService
}

func NewFavoriteController(favoriteService service.FavoriteService, idempotencyService service.IdempotencyService) FavoriteController {
	return &favoriteController{
		favoriteService:    favoriteService,
		idempotencyService: idempotencyService,
	}
}
func (c *favoriteController) Create(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось прочитать данные",
		})
		return
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var createFavorite dto.RequestCreateFavorite
	if err := ctx.ShouldBindJSON(&createFavorite); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось прочитать данные",
		})
		return
	}

	idempotencyKeyRaw := ctx.GetHeader("Idempotency-Key")
	if idempotencyKeyRaw == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не передан ключ идемпотентности",
		})
		return
	}

	idempotencyKey, err := uuid.Parse(idempotencyKeyRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный ключ идемпотентности",
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

	payloadHash := security.CalculatePayloadHash(body)

	idempotency, err := c.idempotencyService.Get(idempotencyKey)
	if err == nil {
		if idempotency.PayloadHash != payloadHash {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "ключ идемпотентности уже использован для другого запроса",
			})
			return
		}

		ctx.Data(
			idempotency.ResponseCode,
			"application/json",
			idempotency.ResponseBody,
		)
		return
	}

	if !errors.Is(err, storage.ErrIdempotencyNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось проверить ключ идемпотентности",
		})
		return
	}

	favorite, err := c.favoriteService.Create(entity.CreateFavorite{
		UserID:   userID,
		FlowerID: createFavorite.FlowerID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseBody, err := json.Marshal(favorite)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось сформировать ответ",
		})
		return
	}

	_, err = c.idempotencyService.Create(entity.CreateIdempotency{
		Key:          idempotencyKey,
		Status:       "completed",
		ResponseCode: http.StatusCreated,
		ResponseBody: responseBody,
		PayloadHash:  payloadHash,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось сохранить ключ идемпотентности",
		})
		return
	}

	ctx.JSON(http.StatusCreated, favorite)
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

	favorites, err := c.favoriteService.GetByUserID(userID)
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

	favorites, err := c.favoriteService.GetByUserIDOld(userID)
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

	favorites, err := c.favoriteService.Delete(favoriteID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, favorites)
}
