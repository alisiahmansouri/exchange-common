package model

import "errors"

var (
	ErrInternal = errors.New("خطای داخلی سرور")
)

// --- خطاهای ثبت‌نام کاربر ---
var (
	ErrInvalidIdentifier    = errors.New("شناسه ایمیل یا شماره موبایل نامعتبر است")
	ErrEmailExists          = errors.New("ایمیل وارد شده قبلاً ثبت شده است")
	ErrPhoneExists          = errors.New("شماره موبایل وارد شده قبلاً ثبت شده است")
	ErrEmailOrPhoneRequired = errors.New("ایمیل یا شماره موبایل الزامی است")
	ErrInvalidEmailFormat   = errors.New("فرمت ایمیل نامعتبر است")
	ErrInvalidPhone         = errors.New("فرمت شماره موبایل نامعتبر است")
	ErrPasswordHash         = errors.New("خطا در رمزنگاری رمز عبور")
	ErrUserCreate           = errors.New("خطا در ایجاد کاربر")
)

// --- خطاهای مربوط به کد تایید (Verification Code) ---
var (
	ErrVerificationCodeInvalidOrExpired = errors.New("کد تایید نامعتبر یا منقضی شده است")
	ErrVerificationCodeAlreadyUsed      = errors.New("کد تایید قبلا استفاده شده است")
	ErrVerificationCodeExpired          = errors.New("کد تایید منقضی شده است")
	ErrVerificationCodeInvalid          = errors.New("کد تایید نامعتبر است")
	ErrVerificationCodeNotFound         = errors.New("کد تایید یافت نشد")
	ErrVerificationCodeNotGenerated     = errors.New("کد تایید تولید نشد")
	ErrVerificationCodeHashMismatch     = errors.New("کد وارد شده صحیح نیست")
	ErrVerificationPurposeInvalid       = errors.New("هدف کد تایید نامعتبر است")
	ErrVerificationChannelInvalid       = errors.New("کانال ارسال کد تایید نامعتبر است")
	ErrVerificationIdentifierInvalid    = errors.New("ایمیل یا موبایل نامعتبر است")
	ErrVerificationRateLimitExceeded    = errors.New("تعداد ارسال کد بیش از حد مجاز است")
	ErrTooManyAttempts                  = errors.New("تعداد تلاش بیش از حد مجاز است")
)

// --- خطاهای مرتبط با ریست پسورد ---
var (
	ErrResetPasswordInvalidOrExpiredCode = errors.New("کد بازیابی رمزعبور نامعتبر یا منقضی شده است")
	ErrResetPasswordCodeUsed             = errors.New("کد بازیابی رمزعبور قبلا استفاده شده است")
	ErrUserNotFound                      = errors.New("کاربر یافت نشد")
)

// --- خطاهای کاربر و 2FA ---
var (
	ErrUserExists   = errors.New("کاربر با این ایمیل قبلا ثبت شده است")
	ErrInvalidCreds = errors.New("اطلاعات ورود نادرست است")
	ErrUserInactive = errors.New("کاربر غیر فعال است")

	Err2FACodeInvalid = errors.New("کد 2FA معتبر نیست")
	Err2FACodeUsed    = errors.New("کد 2FA قبلا استفاده شده است")
	Err2FACodeExpired = errors.New("کد 2FA منقضی شده است")
)

var (
	ErrDepositAmountInvalid    = errors.New("مبلغ واریز باید بزرگتر از صفر باشد")
	ErrWithdrawAmountInvalid   = errors.New("مبلغ برداشت باید بزرگتر از صفر باشد")
	ErrWalletNotFound          = errors.New("کیف پول یافت نشد")
	ErrWalletUnauthorized      = errors.New("دسترسی به کیف پول غیرمجاز است")
	ErrInsufficientFunds       = errors.New("موجودی کافی نیست")
	ErrWalletInactive          = errors.New("کیف پول غیر فعال است")
	ErrAmountInvalid           = errors.New("مقدار نامعتبر است")
	ErrFrozenInsufficientFunds = errors.New("موجودی فریز شده کافی نیست")
	ErrWalletAlreadyExists     = errors.New("این کیف پول قبلا ایجاد شده است")
	ErrBulkOpInvalidType       = errors.New("نوع عملیات گروهی نامعتبر است")
)
