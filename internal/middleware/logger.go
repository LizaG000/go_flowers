package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		slog.Info(
			"HTTP request completed",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.Int("status", ctx.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.String("client_ip", ctx.ClientIP()),
		)
	}
}
