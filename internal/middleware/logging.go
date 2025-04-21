package middleware

import (
	"ienergy-template-go/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logger.StandardLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		entry := logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		c.Next()

		duration := time.Since(start).Seconds()
		entry.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"duration": duration,
		}).Info("request handled")
	}
}
