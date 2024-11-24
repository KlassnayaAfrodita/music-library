package middleware

import (
	"time"

	"music-library/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Логирование входящего запроса
		logger.Log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"query":  c.Request.URL.RawQuery,
			"ip":     c.ClientIP(),
		}).Info("Incoming request")

		// Выполнение обработчика
		c.Next()

		// Логирование ответа
		duration := time.Since(startTime)
		logger.Log.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"duration": duration.Milliseconds(),
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
		}).Info("Response sent")
	}
}
