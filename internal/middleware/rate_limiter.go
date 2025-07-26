package middleware

import (
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"fmt"
)

type RateLimiter struct {
	redisClient *redis.Client
	limit       int           // تعداد درخواست مجاز
	window      time.Duration // بازه زمانی برای limit
}

func NewRateLimiter(redisClient *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		limit:       limit,
		window:      window,
	}
}

func (r *RateLimiter) LimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		count, err := r.redisClient.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "خطای سرور"})
			return
		}

		if count >= r.limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "تعداد درخواست‌ها زیاد است، لطفا بعدا تلاش کنید.",
			})
			return
		}

		if err == redis.Nil {
			// کلید وجود نداره → مقدار 1 با TTL تنظیم کن
			err = r.redisClient.Set(ctx, key, 1, r.window).Err()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "خطای سرور"})
				return
			}
		} else {
			// کلید وجود داره → فقط افزایش بده
			err = r.redisClient.Incr(ctx, key).Err()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "خطای سرور"})
				return
			}
		}

		c.Next()
	}
}
