package captcha

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// Store defines interface for storing captcha codes.
// (برای تست‌نویسی راحت، می‌تونی اینترفیس جدا تعریف کنی)
type Store interface {
	Set(ctx context.Context, requestID string, value string) error
	SetWithTTL(ctx context.Context, requestID string, value string, ttl time.Duration) error
	Get(ctx context.Context, requestID string) (string, error)
	Delete(ctx context.Context, requestID string) error
}

// CaptchaStore implements Store interface using Redis.
type CaptchaStore struct {
	redisClient *redis.Client
	expire      time.Duration
}

// NewCaptchaStore ساخت اینستنس جدید با Redis و مدت TTL پیش‌فرض.
func NewCaptchaStore(redisClient *redis.Client, expire time.Duration) *CaptchaStore {
	return &CaptchaStore{
		redisClient: redisClient,
		expire:      expire,
	}
}

// Set ذخیره مقدار با TTL پیش‌فرض.
func (s *CaptchaStore) Set(ctx context.Context, requestID string, value string) error {
	key := s.redisKey(requestID)
	return s.redisClient.Set(ctx, key, value, s.expire).Err()
}

// SetWithTTL ذخیره مقدار با TTL سفارشی.
func (s *CaptchaStore) SetWithTTL(ctx context.Context, requestID string, value string, ttl time.Duration) error {
	key := s.redisKey(requestID)
	return s.redisClient.Set(ctx, key, value, ttl).Err()
}

// Get مقدار کپچا را می‌خواند.
func (s *CaptchaStore) Get(ctx context.Context, requestID string) (string, error) {
	key := s.redisKey(requestID)
	return s.redisClient.Get(ctx, key).Result()
}

// Delete مقدار کپچا را حذف می‌کند.
func (s *CaptchaStore) Delete(ctx context.Context, requestID string) error {
	key := s.redisKey(requestID)
	return s.redisClient.Del(ctx, key).Err()
}

func (s *CaptchaStore) redisKey(requestID string) string {
	return fmt.Sprintf("captcha:%s", requestID)
}
