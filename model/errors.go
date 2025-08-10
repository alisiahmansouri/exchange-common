package model

import "errors"

// --- خطای عمومی ---
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
	ErrVerificationCodeAlreadyUsed      = errors.New("کد تایید قبلاً استفاده شده است")
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
	ErrResetPasswordCodeUsed             = errors.New("کد بازیابی رمزعبور قبلاً استفاده شده است")
	ErrUserNotFound                      = errors.New("کاربر یافت نشد")
)

// --- خطاهای کاربر و 2FA ---
var (
	ErrUserExists   = errors.New("کاربر با این ایمیل قبلاً ثبت شده است")
	ErrInvalidCreds = errors.New("اطلاعات ورود نامعتبر است")
	ErrUserInactive = errors.New("کاربر غیر فعال است")

	Err2FACodeInvalid = errors.New("کد ورود دو مرحله‌ای نامعتبر است")
	Err2FACodeUsed    = errors.New("کد ورود دو مرحله‌ای قبلاً استفاده شده است")
	Err2FACodeExpired = errors.New("کد ورود دو مرحله‌ای منقضی شده است")
)

// --- خطاهای کیف پول ---
var (
	ErrDepositAmountInvalid     = errors.New("مبلغ واریز باید بزرگتر از صفر باشد")
	ErrWithdrawAmountInvalid    = errors.New("مبلغ برداشت باید بزرگتر از صفر باشد")
	ErrWalletNotFound           = errors.New("کیف پول یافت نشد")
	ErrWalletUnauthorized       = errors.New("دسترسی به کیف پول غیرمجاز است")
	ErrInsufficientFunds        = errors.New("موجودی کافی نیست")
	ErrWalletInactive           = errors.New("کیف پول غیر فعال است")
	ErrAmountInvalid            = errors.New("مقدار وارد شده نامعتبر است")
	ErrFrozenInsufficientFunds  = errors.New("موجودی فریز شده کافی نیست")
	ErrWalletAlreadyExists      = errors.New("این کیف پول قبلاً ایجاد شده است")
	ErrBulkOpInvalidType        = errors.New("نوع عملیات گروهی نامعتبر است")
	ErrWalletNotFoundOrInactive = errors.New("کیف پول یافت نشد یا فعال نیست")
)

// --- خطاهای سفارش و جفت‌ارز ---
var (
	ErrOrderAmountInvalid            = errors.New("مقدار سفارش نامعتبر است")
	ErrOrderLimitPriceRequired       = errors.New("قیمت سفارش limit باید مشخص باشد")
	ErrOrderUserIDInvalid            = errors.New("شناسه کاربر سفارش نامعتبر است")
	ErrPairIDInvalid                 = errors.New("شناسه جفت ارز نامعتبر است")
	ErrOrderSideInvalid              = errors.New("نوع سفارش (خرید یا فروش) نامعتبر است")
	ErrPairNotFoundOrInactive        = errors.New("جفت ارز یافت نشد یا فعال نیست")
	ErrOrderAmountOutOfRange         = errors.New("مقدار سفارش خارج از بازه مجاز است")
	ErrOrderNotFound                 = errors.New("سفارش یافت نشد")
	ErrOrderCannotBeCanceled         = errors.New("این سفارش امکان لغو شدن ندارد")
	ErrOrderTypeInvalid              = errors.New("نوع سفارش نامعتبر است")
	ErrOrderPriceNotAllowedForMarket = errors.New("ارسال قیمت برای سفارش market مجاز نیست")
)

// --- خطاهای سفارش (Order Errors) ---
var (
	ErrOrderInsufficientFunds    = errors.New("موجودی کافی برای ثبت سفارش وجود ندارد")
	ErrOrderConflict             = errors.New("سفارش متناقض یا تکراری است")
	ErrOrderTooManyRequests      = errors.New("تعداد درخواست‌های سفارش بیش از حد مجاز است")
	ErrOrderTimeout              = errors.New("ثبت سفارش به علت محدودیت زمانی انجام نشد")
	ErrOrderCreate               = errors.New("خطا در ثبت سفارش")
	ErrOrderStatusInvalid        = errors.New("وضعیت سفارش نامعتبر است")
	ErrOrderAlreadyCanceled      = errors.New("سفارش قبلاً لغو شده است")
	ErrOrderAlreadyCompleted     = errors.New("سفارش قبلاً تکمیل شده است")
	ErrOrderFillAmountInvalid    = errors.New("مقدار پر شده سفارش نامعتبر است")
	ErrOrderUnfreezingFailed     = errors.New("آزادسازی مبلغ بلوکه شده با خطا مواجه شد")
	ErrOrderWalletFreezeFailed   = errors.New("فریز کردن مبلغ سفارش با خطا مواجه شد")
	ErrOrderMaxOpenOrdersReached = errors.New("حداکثر تعداد سفارش باز مجاز برای کاربر پر شده است")
)
