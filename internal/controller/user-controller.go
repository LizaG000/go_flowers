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
	user, err := c.userService.GetByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}
