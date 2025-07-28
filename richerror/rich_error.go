package richerror

import (
	"errors"
	"fmt"
	"runtime"
)

// Kind نوع معنایی خطا را مشخص می‌کند
type Kind string

const (
	KindValidation      Kind = "validation"
	KindNotFound        Kind = "not_found"
	KindUnauthorized    Kind = "unauthorized"
	KindInternal        Kind = "internal"
	KindForbidden       Kind = "forbidden"
	KindInvalid         Kind = "invalid"
	KindConflict        Kind = "conflict"
	KindTooManyRequests Kind = "too_many_requests"
	KindTimeout         Kind = "timeout"
)

// RichError ساختار کامل برای حمل خطا با زمینه‌های مفید
type RichError struct {
	Op          string // عملیات یا تابع ایجادکننده خطا
	UserMessage string // پیام کاربرپسند
	Code        string // کد یکتا برای ارسال به کلاینت
	Kind        Kind   // نوع خطا (مثلاً validation)
	Err         error  // خطای اصلی
	Caller      string // موقعیت دقیق در کد (مثلاً file.go:42)
	Msg         string // توضیح فنی برای لاگ
}

// Error رشته متنی خطا را برمی‌گرداند
func (r *RichError) Error() string {
	if r.Msg != "" {
		return fmt.Sprintf("%s: %s", r.Op, r.Msg)
	}
	if r.Err != nil {
		return fmt.Sprintf("%s: %v", r.Op, r.Err)
	}
	return r.Op
}

// Unwrap پشتیبانی از errors.Unwrap
func (r *RichError) Unwrap() error {
	return r.Err
}

// Is پشتیبانی از errors.Is
func (r *RichError) Is(target error) bool {
	return errors.Is(r.Err, target)
}

// As پشتیبانی از errors.As
func (r *RichError) As(target interface{}) bool {
	return errors.As(r.Err, target)
}

// String بازنمایی کامل برای debug/log
func (r *RichError) String() string {
	return fmt.Sprintf("[Kind=%s, Op=%s, HashedCode=%s, Msg=%s, UserMsg=%s, Caller=%s]", r.Kind, r.Op, r.Code, r.Msg, r.UserMessage, r.Caller)
}

// New ایجاد RichError جدید از ابتدا
func New(op, userMsg, code string, kind Kind, err error) *RichError {
	return &RichError{
		Op:          op,
		UserMessage: userMsg,
		Code:        code,
		Kind:        kind,
		Err:         err,
		Msg:         extractMsg(err),
		Caller:      callerInfo(),
	}
}

// Wrap یک خطای موجود را به RichError تبدیل می‌کند
func Wrap(op string, err error, userMsg, code string, kind Kind) *RichError {
	if err == nil {
		return nil
	}
	return &RichError{
		Op:          op,
		Err:         err,
		UserMessage: userMsg,
		Code:        code,
		Kind:        kind,
		Msg:         err.Error(),
		Caller:      callerInfo(),
	}
}

// ابزارهای داخلی

func extractMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func callerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", file, line)
}
