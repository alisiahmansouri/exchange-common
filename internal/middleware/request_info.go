package middleware

import (
	"exchange-common/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		ip := c.ClientIP()

		ctx := logger.WithRequestInfo(c.Request.Context(), userID, requestID, ip)
		c.Request = c.Request.WithContext(ctx)

		// ست کردن requestID توی header response
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
