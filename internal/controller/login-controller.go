package controller

import (
	"net/http"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/dto"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/entity"
	"gilab.com/pragmaticrewies/golang-gin-poc/internal/service"

	"github.com/gin-gonic/gin"
)

type LoginController interface {
	Login(ctx *gin.Context)
	Registration(ctx *gin.Context)
}

type loginController struct {
	loginService service.LoginService
}

func NewLoginController(loginService service.LoginService) LoginController {
	return &loginController{
		loginService: loginService,
	}
}

// Login godoc
// @Summary Вход в аккаунт
// @Description Выполняет авторизацию пользователя по email и паролю и возвращает JWT-токен.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Данные для входа"
// @Success 200 {string} string "JWT access token"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (c *loginController) Login(ctx *gin.Context) {
	var login dto.LoginRequest

	if err := ctx.ShouldBindJSON(&login); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось прочитать данные",
		})
		return
	}

	token, err := c.loginService.Login(login.Email, login.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "неверный email или пароль",
		})
		return
	}

	ctx.JSON(http.StatusOK, token)
}

// Registration godoc
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя и возвращает JWT-токен для дальнейшей работы с API.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RequestCreateUser true "Данные нового пользователя"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/registration [post]
func (c *loginController) Registration(ctx *gin.Context) {
	var createUser dto.RequestCreateUser

	if err := ctx.ShouldBindJSON(&createUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось прочитать данные",
		})
		return
	}

	user, err := c.loginService.Registration(entity.CreateUser{
		FirstName:  createUser.FirstName,
		SecondName: createUser.SecondName,
		LastName:   createUser.LastName,
		Email:      createUser.Email,
		BirthDate:  createUser.BirthDate,
	}, createUser.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := c.loginService.Login(
		createUser.Email,
		createUser.Password,
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "неверный email или пароль",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user":         user,
		"access_token": token,
	})
}
