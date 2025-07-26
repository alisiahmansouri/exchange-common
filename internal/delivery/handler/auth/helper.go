package auth

import (
	"crypto/sha256"
	"exchange-common/internal/captcha"
	"exchange-common/internal/consts"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"

	"context"
	"exchange-common/internal/richerror"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

func (h *Handler) verifyCaptchaWrapper(c *gin.Context, op string, ctx context.Context, store *captcha.CaptchaStore, captchaID, captchaAns string) bool {
	if err := h.verifyCaptcha(ctx, store, captchaID, captchaAns, op); err != nil {
		richerror.HandleWrap(c, op, consts.ErrCaptchaInvalid, consts.CodeInvalidCaptcha, richerror.KindInvalid, err)
		return true
	}
	return false
}

func (h *Handler) verifyCaptcha(ctx context.Context, store *captcha.CaptchaStore, captchaID, captchaAns, op string) error {
	if captchaID == "" || captchaAns == "" {
		return richerror.New(op, consts.ErrCaptchaInvalidBody, consts.CodeCaptchaEmpty, richerror.KindInvalid, nil)
	}

	expected, err := store.Get(ctx, captchaID)
	if err != nil {
		return richerror.New(op, consts.ErrCaptchaNotFound, consts.CodeCaptchaNotFound, richerror.KindInvalid, err)
	}

	if strings.ToLower(expected) != strings.ToLower(captchaAns) {
		return richerror.New(op, consts.ErrCaptchaInvalid, consts.CodeCaptchaWrong, richerror.KindUnauthorized, nil)
	}

	_ = store.Delete(ctx, captchaID)
	return nil
}

func (h *Handler) generateTokens(c *gin.Context, op string, userID uuid.UUID) (access, refresh string, hasErr bool) {
	access, refresh, err := h.jwtService.GenerateToken(userID)
	if err != nil {
		richerror.HandleWrap(c, op, consts.ErrAuthTokenGenFail, consts.CodeTokenGenFail, richerror.KindInternal, err)
		return "", "", true
	}
	return access, refresh, false
}

func (h *Handler) canSendVerification(ctx context.Context, identifier string) (bool, error) {
	const (
		maxAttempts    = 5
		expireDuration = 10 * time.Minute
		prefix         = "verification_rate_limit"
	)

	// هش کردن identifier (با حروف کوچک) برای امنیت بیشتر
	hash := sha256.Sum256([]byte(strings.ToLower(identifier)))
	key := fmt.Sprintf("%s:%x", prefix, hash)

	// گرفتن تعداد ارسال فعلی
	countStr, err := h.redisClient.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("redis get failed: %w", err)
	}

	var count int
	if countStr != "" {
		count, _ = strconv.Atoi(countStr) // اگر parse نشه، صفر می‌مونه
	}
	if count >= maxAttempts {
		return false, nil
	}

	// افزایش شمارش و ست کردن مدت انقضا
	pipe := h.redisClient.TxPipeline()
	pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, expireDuration)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("redis pipeline exec failed: %w", err)
	}

	return true, nil
}
