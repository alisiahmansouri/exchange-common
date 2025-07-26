package handler

import (
	"exchange-common/internal/captcha"
	"exchange-common/internal/middleware"
	"exchange-common/internal/service"
	"exchange-common/internal/service/jwt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

type Handler struct {
	svc          *service.Service
	jwtService   *jwt.Service
	captchaStore *captcha.CaptchaStore
	redisClient  *redis.Client
}

func SetupRoutes(router *gin.RouterGroup, svc *service.Service, jwtService *jwt.Service, captchaStore *captcha.CaptchaStore, redisClient *redis.Client) {
	h := &Handler{
		svc:          svc,
		jwtService:   jwtService,
		captchaStore: captchaStore,
		redisClient:  redisClient,
	}

	// اعمال middleware کلی (مثلاً timeout عمومی برای همه گروه‌ها)
	router.Use(middleware.TimeoutMiddleware(5 * time.Second))

	h.registerCaptchaRoutes(router)
	h.registerCurrencyRoutes(router)
	h.registerWalletRoutes(router)
}

func (h *Handler) registerCaptchaRoutes(router *gin.RouterGroup) {
	rateLimiter := middleware.NewRateLimiter(h.redisClient, 5, time.Minute)

	captchaGroup := router.Group("/captcha").
		Use(middleware.TimeoutMiddleware(3 * time.Second)).
		Use(middleware.RequestInfoMiddleware()).
		Use(middleware.RequestLogger()).
		Use(rateLimiter.LimitMiddleware())

	captchaGroup.GET("", h.CaptchaGet)
	captchaGroup.POST("/verify", h.CaptchaVerify)
}

func (h *Handler) registerCurrencyRoutes(router *gin.RouterGroup) {
	currency := router.Group("/currencies").
		Use(middleware.AuthMiddleware(h.jwtService)).
		Use(middleware.RequestInfoMiddleware()).
		Use(middleware.RequestLogger())

	currency.GET("", h.ListCurrencies)
}

func (h *Handler) registerWalletRoutes(router *gin.RouterGroup) {
	wallet := router.Group("/wallets").
		Use(middleware.AuthMiddleware(h.jwtService)).
		Use(middleware.RequestInfoMiddleware()).
		Use(middleware.RequestLogger())

	wallet.GET("", h.ListWallets)
	wallet.POST("/:wallet_id/deposit", h.DepositWallet)
	wallet.POST("/:wallet_id/withdraw", h.WithdrawWallet)
}
