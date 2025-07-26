package richerror

import (
	"errors"
	"exchange-common/internal/logger"
	"exchange-common/internal/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Handle یک خطای جدید را لاگ و به کلاینت ارسال می‌کند
func Handle(c *gin.Context, op, userMsg, code string, kind Kind, err error) {
	richErr := New(op, userMsg, code, kind, err)
	logRichError(c, richErr)
	HTTPErrorHandler(c, richErr)
}

// HandleWrap خطا را wrap کرده، لاگ و ارسال می‌کند
func HandleWrap(c *gin.Context, op, userMsg, code string, kind Kind, err error) {
	richErr := Wrap(op, err, userMsg, code, kind)
	logRichError(c, richErr)
	HTTPErrorHandler(c, richErr)
}

// logRichError لاگ‌گیری ساختاریافته‌ی خطای غنی‌شده
func logRichError(c *gin.Context, richErr *RichError) {
	log := logger.FromContext(c.Request.Context())
	log.Error("❌ Operation failed",
		zap.Error(richErr.Err),
		zap.String("rich_error_code", richErr.Code),
		zap.String("rich_error_msg", richErr.UserMessage),
		zap.String("rich_error_op", richErr.Op),
		zap.String("rich_error_kind", string(richErr.Kind)),
		zap.String("rich_error_caller", richErr.Caller),
	)
}

// HTTPErrorHandler پاسخ JSON مناسب را به کلاینت ارسال می‌کند
func HTTPErrorHandler(c *gin.Context, err error) {
	re := As(err)
	if re == nil {
		model.ErrorResponse(c, http.StatusInternalServerError, "خطای ناشناخته", "INTERNAL_UNKNOWN")
		return
	}

	status, msg, code := ToHTTPResponse(re)
	model.ErrorResponse(c, status, msg, code)
}

// ToHTTPResponse تبدیل RichError به وضعیت HTTP و پیام کاربرپسند
func ToHTTPResponse(err *RichError) (status int, message string, code string) {
	switch err.Kind {
	case KindInvalid, KindValidation:
		return http.StatusBadRequest, err.UserMessage, err.Code
	case KindNotFound:
		return http.StatusNotFound, err.UserMessage, err.Code
	case KindUnauthorized:
		return http.StatusUnauthorized, err.UserMessage, err.Code
	case KindForbidden:
		return http.StatusForbidden, err.UserMessage, err.Code
	case KindConflict:
		return http.StatusConflict, err.UserMessage, err.Code
	case KindTooManyRequests:
		return http.StatusTooManyRequests, err.UserMessage, err.Code
	default:
		return http.StatusInternalServerError, err.UserMessage, err.Code
	}
}

// As تبدیل هر error به RichError (در صورت امکان)
func As(err error) *RichError {
	var re *RichError
	if errors.As(err, &re) {
		return re
	}
	return nil
}
