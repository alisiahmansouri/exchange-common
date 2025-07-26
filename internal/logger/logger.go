package logger

import (
	"context"
	"exchange-common/config"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ctxKey struct{}

var globalLogger *zap.Logger

// Init sets up global zap.Logger with config
func Init(cfg config.Log) {
	if !cfg.Enable {
		globalLogger = zap.NewNop()
		return
	}

	level := parseLevel(cfg.Level)

	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder(cfg.TimeLayout),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.FileMaxSize,
		MaxBackups: cfg.FileMaxBackups,
		MaxAge:     cfg.FileMaxAge,
		Compress:   cfg.FileCompress,
	})
	consoleWriter := zapcore.Lock(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, level),
		zapcore.NewCore(encoder, consoleWriter, level),
	)

	opts := []zap.Option{zap.AddCaller()}
	if cfg.Trace {
		opts = append(opts, zap.AddCallerSkip(1))
	}

	globalLogger = zap.New(core, opts...)
}

// Inject returns new context with zap.Logger
func Inject(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, globalLogger)
}

// FromContext returns zap.Logger with optional user/request fields
func FromContext(ctx context.Context) *zap.Logger {
	logger := globalLogger
	if ctxLogger, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		logger = ctxLogger
	}

	fields := []zap.Field{}
	if v := ctx.Value("userID"); v != nil {
		fields = append(fields, zap.String("user_id", v.(string)))
	}
	if v := ctx.Value("requestID"); v != nil {
		fields = append(fields, zap.String("request_id", v.(string)))
	}
	if v := ctx.Value("ip"); v != nil {
		fields = append(fields, zap.String("ip", v.(string)))
	}
	return logger.With(fields...)
}

// Helpers
func Info(msg string, fields ...zap.Field) { globalLogger.Info(msg, fields...) }
func Error(err error, fields ...zap.Field) { globalLogger.Error(err.Error(), fields...) }
func Fatal(err error, fields ...zap.Field) { globalLogger.Fatal(err.Error(), fields...) }
func Sync()                                { _ = globalLogger.Sync() }
func Log() *zap.Logger                     { return globalLogger }

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func timeEncoder(layout string) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(layout))
	}
}

func WithRequestInfo(ctx context.Context, userID int64, requestID string, ip string) context.Context {
	ctx = context.WithValue(ctx, "userID", fmt.Sprintf("%d", userID))
	ctx = context.WithValue(ctx, "requestID", requestID)
	ctx = context.WithValue(ctx, "ip", ip)
	return ctx
}
