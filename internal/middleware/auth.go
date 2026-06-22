package middleware

import (
	"net/http"
	"strings"

	"gilab.com/pragmaticrewies/golang-gin-poc/internal/security"

	"github.com/gin-gonic/gin"
)

const userIDContextKey = "user_id"

func Auth(publicKeyPath string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		if header == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "токен авторизации не передан",
			})
			ctx.Abort()
			return
		}

		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(header, bearerPrefix) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "неверный формат токена",
			})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, bearerPrefix)

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "токен авторизации не передан",
			})
			ctx.Abort()
			return
		}

		claims, err := security.ParseToken(tokenString, publicKeyPath)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "некорректный или просроченный токен",
			})
			ctx.Abort()
			return
		}

		ctx.Set(userIDContextKey, claims.UserID)

		ctx.Next()
	}
}
