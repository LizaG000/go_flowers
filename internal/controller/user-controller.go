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
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

// Get godoc
// @Summary Получить профиль текущего пользователя
// @Description Возвращает данные пользователя, определённого по JWT-токену.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} entity.User
// @Failure 401 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [get]
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

	ctx.JSON(http.StatusOK, user)
}
