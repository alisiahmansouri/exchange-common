package delivery

import (
	"exchange-common/internal/captcha"
	"exchange-common/internal/db"
	"exchange-common/internal/delivery/handler"
	"exchange-common/internal/delivery/handler/auth"
	"exchange-common/internal/logger"
	"exchange-common/internal/service"
	"exchange-common/internal/service/auth_service"
	"exchange-common/internal/service/jwt"
	"exchange-common/internal/service/repository"
	"exchange-common/internal/service/verification_service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ExchangeRouterRegister struct {
	engine       *gin.Engine
	db           *db.DB
	jwtService   *jwt.Service
	captchaStore *captcha.CaptchaStore
	redisClient  *redis.Client
}

func NewExchangeRouterRegister(
	engine *gin.Engine,
	db *db.DB,
	jwtService *jwt.Service,
	captchaStore *captcha.CaptchaStore,
	redisClient *redis.Client,
) *ExchangeRouterRegister {
	return &ExchangeRouterRegister{
		engine:       engine,
		db:           db,
		jwtService:   jwtService,
		captchaStore: captchaStore,
		redisClient:  redisClient,
	}
}

func (r *ExchangeRouterRegister) Setup() error {
	logger.Info("📡 Setting up ExchangeRouter...")

	apiGroup := r.engine.Group("/v1")

	repo := repository.NewGormRepository(r.db.DB)

	mainSvc := service.NewService(repo)

	authSvc := auth_service.New(repo)

	verificationSvc := verification_service.New(repo)

	handler.SetupRoutes(apiGroup, mainSvc, r.jwtService, r.captchaStore, r.redisClient)

	auth.SetupAuthRoutes(apiGroup, authSvc, verificationSvc, r.jwtService, r.captchaStore, r.redisClient)

	logger.Info("✅ Exchange HTTP routes registered.")
	return nil
}
