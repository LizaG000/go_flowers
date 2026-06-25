package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController interface {
	Check(ctx *gin.Context)
}

type healthController struct {
	db *sql.DB
}

func NewHealthController(db *sql.DB) HealthController {
	return &healthController{
		db: db,
	}
}

// Check godoc
// @Summary Проверка состояния сервиса
// @Description Внутренняя конечная точка для проверки доступности приложения и базы данных.
// @Tags internal
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /internal/health [get]
func (c *healthController) Check(ctx *gin.Context) {
	if err := c.db.PingContext(ctx.Request.Context()); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "unavailable",
			"database": "unavailable",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": "ok",
	})
}
