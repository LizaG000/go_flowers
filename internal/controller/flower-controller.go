package controller

import (
	"net/http"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/dto"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/service"

	"github.com/gin-gonic/gin"
)

type FlowerController interface {
	GetAll(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type flowerController struct {
	service service.FlowerService
}

func NewFlowerController(service service.FlowerService) FlowerController {
	return &flowerController{
		service: service,
	}
}

// GetAll godoc
// @Summary Получить список цветов
// @Description Возвращает все цветы, доступные в каталоге.
// @Tags flowers
// @Produce json
// @Success 200 {array} entity.Flower
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /flowers [get]
func (c *flowerController) GetAll(ctx *gin.Context) {
	flowers, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, flowers)
}

// Create godoc
// @Summary Создать цветок
// @Description Добавляет новый цветок в каталог магазина.
// @Tags flowers
// @Accept json
// @Produce json
// @Param flower body entity.CreateFlower true "Данные нового цветка"
// @Success 201 {object} entity.Flower
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /flowers [post]
func (c *flowerController) Create(ctx *gin.Context) {
	var createFlower entity.CreateFlower

	if err := ctx.ShouldBindJSON(&createFlower); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось прочитать данные",
		})
		return
	}

	flower, err := c.service.Create(createFlower)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, flower)
}

// Update godoc
// @Summary Обновить цветок
// @Description Обновляет данные существующего цветка.
// @Tags flowers
// @Accept json
// @Produce json
// @Param flower body dto.RequestUpdateFlower true "Новые данные цветка"
// @Success 201 {object} entity.Flower
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /flowers [put]
func (c *flowerController) Update(ctx *gin.Context) {
	var updatedFlower dto.RequestUpdateFlower

	if err := ctx.ShouldBindJSON(&updatedFlower); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось прочитать данные",
		})
		return
	}

	flower, err := c.service.Update(
		updatedFlower.ID,
		updatedFlower.Update_data,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, flower)
}

// Delete godoc
// @Summary Удалить цветок
// @Description Удаляет цветок из каталога по его идентификатору.
// @Tags flowers
// @Accept json
// @Produce json
// @Param flower body entity.DeleteFlower true "Идентификатор удаляемого цветка"
// @Success 200 {object} entity.Flower
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /flowers [delete]
func (c *flowerController) Delete(ctx *gin.Context) {
	var deletedFlower entity.DeleteFlower

	if err := ctx.ShouldBindJSON(&deletedFlower); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось прочитать данные",
		})
		return
	}

	flower, err := c.service.Delete(deletedFlower.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, flower)
}
