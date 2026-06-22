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
func (c *flowerController) Update(ctx *gin.Context) {
	var updatedFlower dto.RequestUpdateFlower

	if err := ctx.ShouldBindJSON(&updatedFlower); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось прочитать данные",
		})
		return
	}

	flower, err := c.service.Update(updatedFlower.ID, updatedFlower.Update_data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, flower)
}

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
	}
	ctx.JSON(http.StatusCreated, flower)
}
