package delivery

import (
	"exchange-common/internal/richerror"
	"github.com/gin-gonic/gin"
)

type ErrorHandler interface {
	Handle(c *gin.Context, op, userMsg, code string, kind richerror.Kind, err error)
	HandleWrap(c *gin.Context, op, userMsg, code string, kind richerror.Kind, err error)
}
