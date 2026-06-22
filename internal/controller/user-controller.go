package controller

import (
	"net/http"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController interface {
	Get(ctx *gin.Context)
}

type userController struct {
	userService  service.UserService
	loginService service.LoginService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (c *userController) Get(ctx *gin.Context) {
	userIDRaw := ctx.Query("userID")

	if userIDRaw == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id пользователя обязателен",
		})
		return
	}
	userID, err := uuid.Parse(userIDRaw)
	user, err := c.userService.GetByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}
