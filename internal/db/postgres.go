package db

import (
	"context"
	"exchange-common/config"
	"exchange-common/internal/entity"
	"exchange-common/internal/logger"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(ctx context.Context, cfg config.Postgres) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Tehran",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	gormZapLogger := logger.NewGormLogger(ctx)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormZapLogger,
	})
	if err != nil {
		logger.FromContext(ctx).Error("❌ failed to connect to postgres", zap.Error(err))
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.FromContext(ctx).Error("❌ failed to get sql.DB from gorm", zap.Error(err))
		return nil, fmt.Errorf("get sql.DB from gorm failed: %w", err)
	}

	// Connection pool configs (می‌تونی این مقادیر رو هم از config بگیری)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.FromContext(ctx).Info("📦 connected to PostgreSQL successfully")

	if err := Migrate(db); err != nil {
		logger.FromContext(ctx).Error("❌ database migration failed", zap.Error(err))
		return nil, fmt.Errorf("database migration failed: %w", err)
	}
	logger.FromContext(ctx).Info("✅ database migrated successfully")

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Currency{},
		&entity.Wallet{},
		&entity.Transaction{},
		&entity.Order{},
		&entity.Pair{},
	)
}
