package health

import (
	"context"
	"exchange-common/internal/model"
	"net/http"
	"time"

	"exchange-common/internal/logger"
	"exchange-common/internal/richerror"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthChecker interface (برای تسهیل تست و توسعه)
type HealthChecker interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	db HealthChecker
}

func NewHealthHandler(db HealthChecker) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck godoc
// @Summary      بررسی سلامت سرویس
// @Description  چک سلامت اتصال به دیتابیس و وضعیت کلی سرور
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.Response
// @Failure      503  {object}  model.Response
// @Router       /health [get]
// @Security     BearerAuth
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	start := time.Now()
	logger.FromContext(ctx).Info("Health check started")

	if err := h.db.Ping(ctx); err != nil {
		logger.FromContext(ctx).Error("Database ping failed", zap.Error(err))

		richErr := richerror.Wrap(
			"HealthHandler.HealthCheck",
			err,
			"اتصال به دیتابیس برقرار نیست",
			"DB_CONN_503",
			richerror.KindInternal,
		)
		richerror.HTTPErrorHandler(c, richErr)
		return
	}

	duration := time.Since(start)

	logger.FromContext(ctx).Info("Health check successful", zap.Duration("duration", duration))

	data := gin.H{
		"status":        "ok",
		"database":      "connected",
		"response_time": duration.String(),
		"timestamp":     time.Now().Format("2006-01-02 15:04:05"),
	}

	model.SuccessResponse(c, http.StatusOK, data, "سرویس سالم است")
}
