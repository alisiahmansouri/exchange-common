package consts

import "time"

// --- Operation Identifiers ---
const (
	OpCaptchaGet    = "Handler.CaptchaGet"
	OpCaptchaVerify = "Handler.CaptchaVerify"
)

// --- Error Messages ---
const (
	ErrCaptchaIDGenFail   = "خطا در تولید شناسه کپچا"
	ErrCaptchaGenFail     = "خطا در تولید کپچا"
	ErrCaptchaStoreFail   = "خطا در ذخیره کپچا"
	ErrCaptchaInvalidBody = "بدنه درخواست نامعتبر است"
	ErrCaptchaInvalidID   = "شناسه کپچا معتبر نیست"
)

// --- Error Codes ---
const (
	CodeCaptchaIDGenerationFailed = "CAPTCHA_ID_GENERATION_FAILED"
	CodeCaptchaGenerationFailed   = "CAPTCHA_GENERATION_FAILED"
	CodeCaptchaStorageFailed      = "CAPTCHA_STORAGE_FAILED"
	CodeInvalidCaptchaID          = "INVALID_CAPTCHA_ID"
)

// --- Success Messages ---
const (
	MsgCaptchaVerified = "کپچا معتبر است"
)

// --- Captcha Config ---
const (
	CaptchaTTL    = 2 * time.Minute
	CaptchaHeight = 60
	CaptchaWidth  = 240
	CaptchaLength = 2
	CaptchaSource = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var CaptchaFonts = []string{
	"wqy-microhei.ttc",
	"comic.ttf",
	"Vazirmatn-Bold.ttf",
}
