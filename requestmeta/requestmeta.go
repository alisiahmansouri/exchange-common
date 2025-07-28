package requestmeta

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type RequestMeta struct {
	Ctx       context.Context
	Logger    *zap.Logger
	StartTime time.Time
}

type key int

const requestMetaKey key = 0

func NewRequestMeta(c *gin.Context) *RequestMeta {
	ctx := c.Request.Context()
	reqID := c.GetString("RequestID")
	if reqID == "" {
		reqID = "unknown" // یا می‌تونی uuid.NewString() بذاری
	}
	log := zap.L().With(zap.String("request_id", reqID))
	start := time.Now()

	meta := &RequestMeta{
		Ctx:       ctx,
		Logger:    log,
		StartTime: start,
	}

	newCtx := context.WithValue(ctx, requestMetaKey, meta)
	c.Request = c.Request.WithContext(newCtx)

	return meta
}

func FromContext(ctx context.Context) *RequestMeta {
	meta, _ := ctx.Value(requestMetaKey).(*RequestMeta)
	return meta
}

func (m *RequestMeta) Elapsed() time.Duration {
	return time.Since(m.StartTime)
}
