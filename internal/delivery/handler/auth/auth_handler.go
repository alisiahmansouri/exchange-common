package auth

import (
	"exchange-common/internal/captcha"
	"exchange-common/internal/middleware"
	"exchange-common/internal/service/auth_service"
	"exchange-common/internal/service/jwt"
	"exchange-common/internal/service/verification_service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

type Handler struct {
	authSvc         *auth_service.Service
	verificationSVC *verification_service.Service
	jwtService      *jwt.Service
	captchaStore    *captcha.CaptchaStore
	redisClient     *redis.Client
}

// راه‌اندازی مسیرهای احراز هویت (Auth)
func SetupAuthRoutes(
	router *gin.RouterGroup,
	asvc *auth_service.Service,
	vSVC *verification_service.Service,
	jwtService *jwt.Service,
	captchaStore *captcha.CaptchaStore,
	redisClient *redis.Client,
) *Handler {
	h := &Handler{
		authSvc:         asvc,
		verificationSVC: vSVC,
		jwtService:      jwtService,
		captchaStore:    captchaStore,
		redisClient:     redisClient,
	}

	router.Use(middleware.TimeoutMiddleware(5 * time.Second))

	auth := router.Group("/auth").
		Use(middleware.RequestInfoMiddleware()).
		Use(middleware.RequestLogger())

	// روت‌های عمومی (بدون نیاز به توکن)
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/forgot-password", h.ForgotPassword) // ⬅️ فراموشی رمز عبور
	auth.POST("/reset-password", h.ResetPassword)   // ⬅️ ریست رمز عبور

	// --- روت‌های تایید جداگانه برای هر نوع ۲FA
	auth.POST("/verify-register-2fa", h.VerifyRegister2FA)
	auth.POST("/verify-login-2fa", h.VerifyLogin2FA)
	auth.POST("/verify-email", h.VerifyEmail2FA)
	auth.POST("/verify-phone", h.VerifyPhone2FA)
	auth.POST("/resend-verification", h.ResendVerification)

	// روت‌های محافظت‌شده (با JWT)
	authProtected := router.Group("/auth").
		Use(middleware.AuthMiddleware(h.jwtService)).
		Use(middleware.RequestInfoMiddleware()).
		Use(middleware.RequestLogger())
	authProtected.POST("/logout", h.Logout)

	return h
}
