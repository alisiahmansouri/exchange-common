package middleware

import (
	"strconv"
	"strings"

	"exchange-common/internal/logger"
	"exchange-common/internal/richerror"
	"exchange-common/internal/service/jwt"

	"github.com/gin-gonic/gin"
)

const (
	ctxKeyUserID = "user_id"
	ctxKeyJWT    = "jwt_token"
	opAuth       = "middleware.AuthMiddleware"
)

func AuthMiddleware(jwtService *jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			handleUnauthorized(c, "هدر Authorization ارسال نشده است", "AUTH_HEADER_MISSING")
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			handleUnauthorized(c, "فرمت هدر Authorization نادرست است (Bearer <token>)", "AUTH_HEADER_FORMAT_INVALID")
			return
		}

		tokenStr := parts[1]

		claims, err := jwtService.ValidateToken(tokenStr)
		if err != nil {
			handleUnauthorizedWrap(c, err, "توکن نامعتبر است یا منقضی شده", "TOKEN_INVALID_OR_EXPIRED")
			return
		}

		// قرار دادن اطلاعات در Gin context
		c.Set(ctxKeyUserID, claims.UserID)
		c.Set(ctxKeyJWT, tokenStr)

		// قرار دادن در context برای logger
		userIDInt, err := strconv.ParseInt(claims.UserID, 10, 64)
		if err != nil {
			handleUnauthorizedWrap(c, err, "فرمت userID نامعتبر است", "INVALID_USER_ID_FORMAT")
			return
		}

		ctx := logger.WithRequestInfo(c.Request.Context(), userIDInt, "", c.ClientIP())
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func handleUnauthorized(c *gin.Context, userMsg, code string) {
	err := richerror.New(
		opAuth,
		userMsg,
		code,
		richerror.KindUnauthorized,
		nil,
	)
	richerror.HTTPErrorHandler(c, err)
	c.Abort()
}

func handleUnauthorizedWrap(c *gin.Context, err error, userMsg, code string) {
	wrapped := richerror.Wrap(
		opAuth,
		err,
		userMsg,
		code,
		richerror.KindUnauthorized,
	)
	richerror.HTTPErrorHandler(c, wrapped)
	c.Abort()
}
