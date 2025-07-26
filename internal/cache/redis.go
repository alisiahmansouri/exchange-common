package cache

import (
	"context"
	"exchange-common/config"
	"exchange-common/internal/logger"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var Client *redis.Client

func InitRedis(ctx context.Context, cfg config.Redis) error {
	if Client != nil {
		_ = Client.Close() // بستن اتصال قبلی در صورت وجود
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := Client.Ping(ctx).Err(); err != nil {
		logger.Log().Error("❌ failed to connect to Redis",
			zap.Error(err),
			zap.String("addr", cfg.Addr),
			zap.Int("db", cfg.DB),
		)
		return fmt.Errorf("redis ping error: %w", err)
	}

	logger.Log().Info("📦 connected to Redis",
		zap.String("addr", cfg.Addr),
		zap.Int("db", cfg.DB),
	)
	return nil
}
