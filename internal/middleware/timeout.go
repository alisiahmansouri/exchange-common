package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"exchange-common/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ساخت کانتکست با تایم‌اوت
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// جایگزینی کانتکست درخواست
		c.Request = c.Request.WithContext(ctx)

		// sync.Once برای جلوگیری از چندبار پاسخ دادن
		var once sync.Once

		// کانال برای اطلاع از پایان پردازش
		finished := make(chan struct{})

		// اجرای هندلر در goroutine جدا
		go func() {
			defer func() {
				// اگر panic شده، اون رو لاگ کن
				if r := recover(); r != nil {
					log := logger.FromContext(ctx)
					log.Error("panic recovered in timeout middleware",
						zap.Any("panic", r),
					)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error":   "internal server error",
						"message": "unexpected panic occurred",
					})
				}
			}()
			c.Next()
			close(finished)
		}()

		select {
		case <-ctx.Done():
			// زمانی‌که تایم‌اوت اتفاق بیفته
			once.Do(func() {
				log := logger.FromContext(ctx)
				log.Warn("⏱️ request timed out",
					zap.String("path", c.FullPath()),
					zap.Duration("timeout", timeout),
				)

				c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
					"error":   "request timeout",
					"message": "Your request took too long to process.",
				})
			})

		case <-finished:
			// هندلر به موقع تموم شد، مشکلی نیست
		}
	}
}
