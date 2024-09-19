package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		// Após a requisição
		duration := time.Since(startTime)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		logger.WithFields(logrus.Fields{
			"status":    status,
			"method":    method,
			"path":      path,
			"duration":  duration,
			"client_ip": clientIP,
			"timestamp": startTime.Format(time.RFC3339),
		}).Info("Request processed")
	}
}
