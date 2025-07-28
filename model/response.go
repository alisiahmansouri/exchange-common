package model

import "github.com/gin-gonic/gin"

type Response struct {
	Code      int         `json:"code"`                 // HTTP status code
	Success   bool        `json:"success"`              // وضعیت موفقیت
	Data      interface{} `json:"data,omitempty"`       // داده (در صورت موفقیت)
	Message   string      `json:"message,omitempty"`    // پیام (موفق یا خطا)
	ErrorCode string      `json:"error_code,omitempty"` // کد یکتا برای خطا (اختیاری)
}

type ErrorResponseStruct struct {
	Code      int    `json:"code"`
	Success   bool   `json:"success"`
	Error     string `json:"error"`      // پیام فنی (developer-friendly)
	Message   string `json:"message"`    // پیام کاربر (user-friendly)
	ErrorCode string `json:"error_code"` // کد خطای اختصاصی (مثلاً validation_error)
}

func SuccessResponse(c *gin.Context, code int, data interface{}, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	res := Response{
		Code:    code,
		Success: true,
		Data:    data,
		Message: message,
	}

	c.JSON(code, res)
}

func ErrorResponse(c *gin.Context, code int, msg, errorCode string) {
	res := Response{
		Code:    code,
		Success: false,
		Message: msg,
	}
	if errorCode != "" {
		// می‌تونی این فیلد رو به ساختار Response اضافه کنی اگر لازم بود
		res.ErrorCode = errorCode
	}
	c.JSON(code, res)
}

type SimpleMessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"کپچا معتبر است"`
}
