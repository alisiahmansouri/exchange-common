// internal/logger/gorm_logger.go
package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type GormZapLogger struct {
	log *zap.Logger
}

func NewGormLogger(ctx context.Context) logger.Interface {
	return &GormZapLogger{
		log: FromContext(ctx),
	}
}

func (l *GormZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	// برای سادگی سطح لاگ رو ذخیره نمی‌کنیم؛ می‌تونی اضافه‌اش کنی
	return l
}

func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.Sugar().Infof(msg, data...)
}

func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.Sugar().Warnf(msg, data...)
}

func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.Sugar().Errorf(msg, data...)
}

func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	if err != nil {
		l.log.Error("gorm error", append(fields, zap.Error(err))...)
		return
	}

	l.log.Debug("gorm query", fields...)
}
