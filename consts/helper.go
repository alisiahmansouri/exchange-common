package consts

import "time"

// consts/auth_handler.go
const (
	MaxLoginAttempts     = 5
	LoginAttemptDuration = 15 * time.Minute
	TwoFADurationMinutes = 5
	TwoFASendChannel     = 2

	Code2FACodeGen         = "2FA_CODE_GEN"
	ErrAuth2FACodeSendFail = "خطا در ارسال کد تایید"
	Code2FASendFail        = "2FA_SEND_FAIL"

	CodeTokenGenFail = "TOKEN_GEN_FAIL"
)

// consts/captcha.go
const (
	ErrCaptchaInvalid  = "کپچا نادرست است"
	CodeInvalidCaptcha = "INVALID_CAPTCHA"

	CodeCaptchaEmpty    = "CAPTCHA_EMPTY"
	ErrCaptchaNotFound  = "کپچا یافت نشد یا منقضی شده است"
	CodeCaptchaNotFound = "CAPTCHA_NOT_FOUND"
	CodeCaptchaWrong    = "CAPTCHA_WRONG"
)
