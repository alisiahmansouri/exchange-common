package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"exchange-common/config"
	"exchange-common/docs"
	"exchange-common/internal/cache"
	"exchange-common/internal/captcha"
	"exchange-common/internal/db"
	"exchange-common/internal/delivery"
	"exchange-common/internal/health"
	"exchange-common/internal/logger"
	"exchange-common/internal/service/jwt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title           Exchange API
// @version         1.0
// @description     Exchange cryptocurrency trading platform
// @host            localhost:7000
// @BasePath        /v1
// @schemes         http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("❌ config load error: %v", err)
	}

	logger.Init(cfg.Log)

	if err := run(cfg); err != nil {
		logger.Fatal(fmt.Errorf("💥 application terminated: %w", err))
	}
}

func run(cfg *config.Config) error {
	ctx := context.Background()
	docs.SwaggerInfo.BasePath = "/v1"

	dbInstance, err := initDatabase(ctx, cfg)
	if err != nil {
		return err
	}

	if err := cache.InitRedis(ctx, cfg.Redis); err != nil {
		return fmt.Errorf("redis init error: %w", err)
	}

	captchaStore := captcha.NewCaptchaStore(cache.Client, 2*time.Minute)

	jwtService := jwt.New(
		cfg.JWT.SecretKey,
		cfg.JWT.TokenExpiration,
		cfg.JWT.RefreshExpiration,
		cache.Client,
	)

	engine := setupHTTPServer(cfg, dbInstance, jwtService, captchaStore)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		logger.Info(fmt.Sprintf("🚀 Server started on port %d", cfg.HTTP.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		logger.Info("🛑 Signal received: " + sig.String())
	case err := <-errChan:
		logger.Error(fmt.Errorf("🔥 server error: %w", err))
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		logger.Error(fmt.Errorf("❌ graceful shutdown failed: %w", err))
	} else {
		logger.Info("✅ Server gracefully stopped")
	}

	return nil
}

func initDatabase(ctx context.Context, cfg *config.Config) (*db.DB, error) {
	gormDB, err := db.New(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("❌ postgres init error: %w", err)
	}

	if err := db.SeedCurrencies(gormDB); err != nil {
		logger.FromContext(ctx).Error("❌ failed to seed currencies", zap.Error(err))
		return nil, err
	}

	logger.FromContext(ctx).Info("✅ currencies seeded successfully")
	return db.NewDB(gormDB), nil
}

func setupHTTPServer(cfg *config.Config, dbInstance *db.DB, jwtService *jwt.Service, captchaStore *captcha.CaptchaStore) *gin.Engine {
	gin.SetMode(cfg.HTTP.Mode)

	router := gin.New()
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")

	healthHandler := health.NewHealthHandler(dbInstance)
	health.RegisterHealthRoutes(v1, healthHandler)

	exchangeRouter := delivery.NewExchangeRouterRegister(router, dbInstance, jwtService, captchaStore, cache.Client)
	if err := exchangeRouter.Setup(); err != nil {
		logger.Fatal(fmt.Errorf("❌ failed to setup exchange router: %w", err))
	}

	if cfg.Env == "development" || cfg.HTTP.Mode == gin.DebugMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}
