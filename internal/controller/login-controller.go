package controller

import (
	"log/slog"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, token)
}

func (c *loginController) Registration(ctx *gin.Context) {
	var createUser dto.RequestCreateUser

	if err := ctx.ShouldBindJSON(&createUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "не удалось прочитать данные",
		})
		return
	}
	slog.Info(
		"входные данные",
		slog.Any("user", createUser),
	)
	user, err := c.loginService.Registration(entity.CreateUser{
		FirstName:  createUser.FirstName,
		SecondName: createUser.SecondName,
		LastName:   createUser.LastName,
		Email:      createUser.Email,
		BirthDate:  createUser.BirthDate,
	}, createUser.Password)
	slog.Info(
		"пользователь создан",
		slog.Any("user", user),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := c.loginService.Login(createUser.Email, createUser.Password)
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
