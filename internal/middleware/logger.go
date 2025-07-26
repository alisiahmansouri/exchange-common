package middleware

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"exchange-common/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := uuid.NewString()
		ip := c.ClientIP()
		userID := c.GetString("user_id") // فرض: از middleware احراز هویت گرفته شده
		userAgent := c.GetHeader("User-Agent")

		// ========== استخراج و پاکسازی بادی درخواست ==========
		body := extractAndSanitizeBody(c)

		// ========== تزریق اطلاعات به context ==========
		ctx := context.WithValue(c.Request.Context(), "requestID", requestID)
		ctx = context.WithValue(ctx, "userID", userID)
		ctx = context.WithValue(ctx, "ip", ip)
		ctx = logger.Inject(ctx)

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		duration := time.Since(start)

		logger.FromContext(ctx).Info("📥 درخواست پردازش شد",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", ip),
			zap.String("user_id", userID),
			zap.String("request_id", requestID),
			zap.String("user_agent", userAgent),
			zap.Duration("latency", duration),
			zap.String("request_body", body),
		)
	}
}

func extractAndSanitizeBody(c *gin.Context) string {
	if c.Request.Body == nil || !strings.Contains(c.GetHeader("Content-Type"), "application/json") {
		return ""
	}

	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(nil))
		return ""
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
	rawBody := string(buf)

	// حذف فیلدهای حساس
	safeBody := sanitizeRequestBody(rawBody)

	// truncate بدنه بزرگ
	const maxBodySize = 1024
	if len(safeBody) > maxBodySize {
		safeBody = safeBody[:maxBodySize] + "...(truncated)"
	}

	return safeBody
}

func sanitizeRequestBody(body string) string {
	sensitiveFields := []string{"password", "token", "captcha_ans", "secret"}

	for _, field := range sensitiveFields {
		re := regexp.MustCompile(`"` + field + `"\s*:\s*"(.*?)"`)
		body = re.ReplaceAllString(body, `"`+field+`":"[REDACTED]"`)
	}

	return body
}
