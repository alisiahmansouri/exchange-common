package health

import (
	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(r gin.IRouter, handler *HealthHandler) {
	r.GET("/health", handler.HealthCheck)
}
