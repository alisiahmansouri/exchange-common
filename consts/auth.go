package consts

// ==== Operation IDs (Swagger + Logging) ====
const (
	OpAuthRegister                = "AuthHandler.Register"
	OpAuthLogin                   = "AuthHandler.Login"
	OpAuthVerify2FA               = "AuthHandler.Verify2FA"         // کلی
	OpAuthVerifyRegister2FA       = "AuthHandler.VerifyRegister2FA" // جدا
	OpAuthVerifyLogin2FA          = "AuthHandler.VerifyLogin2FA"    // جدا
	OpAuthVerifyEmail             = "AuthHandler.VerifyEmail"
	OpAuthVerifyPhone             = "AuthHandler.VerifyPhone"
	OpAuthRefreshToken            = "AuthHandler.RefreshToken"
	OpAuthLogout                  = "AuthHandler.Logout"
	OpAuthForgotPassword          = "AuthHandler.ForgotPassword"
	OpAuthResetPassword           = "AuthHandler.ResetPassword"
	OpAuthResendEmailVerify       = "AuthHandler.ResendEmailVerification"
	OpAuthSendPhoneVerification   = "AuthHandler.SendPhoneVerification"
	OpAuthResendPhoneVerification = "AuthHandler.ResendPhoneVerification"
	OpAuthResendVerification      = "AuthHandler.ResendVerification"
)

// ==== Error Messages (Client & Dev Friendly) ====
const (
	ErrInvalidBody             = "بدنه درخواست نامعتبر است"
	ErrInvalidParams           = "پارامترهای اجباری وارد نشده یا نامعتبر است"
	ErrInvalidEmail            = "فرمت ایمیل نامعتبر است"
	ErrInvalidPassword         = "رمز جدید معتبر نیست"
	ErrInvalidOrExpiredCode    = "کد بازیابی اشتباه یا منقضی شده"
	ErrResetPasswordFail       = "ریست رمز عبور ناموفق"
	ErrAuthTokenGenFail        = "خطا در تولید توکن"
	ErrAuthUserExists          = "کاربر با این ایمیل یا موبایل قبلاً ثبت شده است"
	ErrTooManyAttempts         = "تعداد تلاش‌های ناموفق بیش از حد مجاز است"
	ErrAuthInvalidCredentials  = "ایمیل، موبایل یا رمز عبور اشتباه است"
	ErrAuthInvalidUserID       = "شناسه کاربر نامعتبر است"
	ErrAuthInvalidRefreshToken = "توکن refresh نامعتبر یا منقضی شده است"
	ErrAuthEmptyToken          = "توکن ارسال نشده یا خالی است"
	ErrAuthNoAuthHeader        = "هدر Authorization ارسال نشده است"
	ErrAuthUserIDTypeInvalid   = "نوع شناسه کاربر نامعتبر است"
	ErrAuthTokenNotFound       = "توکن در context یافت نشد"
	ErrAuthRevokeFail          = "خطا در لغو توکن"
	Err2FACodeCheckFail        = "خطا در بررسی وضعیت تأیید دو عاملی"
	Err2FACodeInvalid          = "کد تایید دو عاملی خالی یا نامعتبر است"
	ErrLoginThrottled          = "محدودیت ورود به دلیل تلاش بیش از حد"
	Err2FASendFail             = "ارسال کد تایید دو عاملی با خطا مواجه شد"
	ErrVerificationUsed        = "کد تایید قبلا استفاده شده است"
	ErrInvalidPurpose          = "هدف ارسال کد تایید نامعتبر است"
	ErrInvalidChannel          = "کانال ارسال کد تایید نامعتبر است"
	ErrRateLimitCheckFail      = "خطا در بررسی محدودیت ارسال"
	ErrRateLimitExceeded       = "تعداد ارسال کد بیش از حد مجاز است، لطفا بعدا تلاش کنید"
	ErrFindUserFail            = "خطا در پیدا کردن کاربر"
	ErrUserNotFound            = "کاربر با این شناسه یافت نشد"
	ErrEmailNotVerified        = "ایمیل شما هنوز تایید نشده است"
	ErrPhoneNotVerified        = "شماره موبایل شما هنوز تایید نشده است"
	// Forgot/Reset
	ErrForgotPasswordEmailNotFound = "ایمیل وارد شده یافت نشد"
	ErrForgotPasswordSendFail      = "خطا در ارسال ایمیل بازیابی"
	// Phone
	ErrPhoneInvalid                = "شماره موبایل نامعتبر"
	ErrSendPhoneVerificationFail   = "ارسال کد تایید موبایل ناموفق"
	ErrVerifyPhoneFail             = "تایید شماره موبایل ناموفق"
	ErrResendPhoneVerificationFail = "ارسال مجدد کد تایید موبایل ناموفق"
	// Email Verification
	ErrVerifyEmailFail             = "تایید ایمیل ناموفق"
	ErrResendEmailVerificationFail = "ارسال مجدد ایمیل تایید ناموفق"
	ErrInvalidEmailOrPhone         = "فرمت ایمیل یا موبایل نامعتبر است"
	ErrEmailExists                 = "ایمیل وارد شده قبلاً ثبت شده است"
	ErrPhoneExists                 = "شماره موبایل وارد شده قبلاً ثبت شده است"
	Err2FAAttemptCheckFail         = "خطا در بررسی محدودیت تلاش دو عاملی"
)

// ==== Error Messages Aliases for Compatibility ====
const (
	ErrAuthInvalidBody    = ErrInvalidBody
	ErrAuthInvalidParams  = ErrInvalidParams
	ErrAuthInvalidEmail   = ErrInvalidEmail
	ErrAuthInvalidPhone   = ErrPhoneInvalid
	ErrAuthInvalidCode    = ErrInvalidOrExpiredCode
	ErrAuthInvalidCaptcha = "کپچای وارد شده نامعتبر است"
	ErrAuth2FACodeInvalid = Err2FACodeInvalid
)

// ==== Error Codes (Log & Client) ====
const (
	CodeInvalidBody            = "INVALID_BODY"
	CodeInvalidParams          = "INVALID_PARAMS"
	CodeInvalidEmail           = "INVALID_EMAIL"
	CodeInvalidPassword        = "INVALID_PASSWORD"
	CodeInvalidOrExpiredCode   = "INVALID_OR_EXPIRED_CODE"
	CodeResetPasswordFail      = "RESET_PASSWORD_FAIL"
	CodeAuthTokenGenFail       = "AUTH_TOKEN_GEN_FAIL"
	CodeUserAlreadyExists      = "USER_ALREADY_EXISTS"
	CodeInvalidCredentials     = "INVALID_CREDENTIALS"
	CodeTooManyAttempts        = "TOO_MANY_ATTEMPTS"
	CodeInvalidUserID          = "INVALID_USER_ID"
	CodeInvalidRefreshToken    = "INVALID_REFRESH_TOKEN"
	CodeTokenEmpty             = "TOKEN_EMPTY"
	CodeAuthHeaderMissing      = "AUTH_HEADER_MISSING"
	CodeUserIDTypeInvalid      = "USER_ID_TYPE_INVALID"
	CodeTokenNotFoundInContext = "TOKEN_NOT_FOUND_IN_CONTEXT"
	CodeJWTRevokeError         = "JWT_REVOKE_ERROR"
	Code2FACheckError          = "2FA_CHECK_ERROR"
	CodeInvalid2FACode         = "INVALID_2FA_CODE"
	CodeLoginThrottled         = "LOGIN_THROTTLED"
	CodeInvalidPurpose         = "INVALID_PURPOSE"
	CodeInvalidChannel         = "INVALID_CHANNEL"
	CodeRateLimitCheckFail     = "RATE_LIMIT_CHECK_FAIL"
	CodeRateLimitExceeded      = "RATE_LIMIT_EXCEEDED"
	Code2FAAttemptCheckFail    = "2FA_ATTEMPT_CHECK_FAIL"
	// Forgot/Reset
	CodeForgotPasswordEmailNotFound = "FORGOT_PASSWORD_EMAIL_NOT_FOUND"
	CodeForgotPasswordSendFail      = "FORGOT_PASSWORD_SEND_FAIL"
	// Phone
	CodeInvalidPhone                = "INVALID_PHONE"
	CodeSendPhoneVerificationFail   = "SEND_PHONE_VERIFICATION_FAIL"
	CodeVerifyPhoneFail             = "VERIFY_PHONE_FAIL"
	CodeResendPhoneVerificationFail = "RESEND_PHONE_VERIFICATION_FAIL"
	// Email Verification
	CodeVerifyEmailFail             = "VERIFY_EMAIL_FAIL"
	CodeResendEmailVerificationFail = "RESEND_EMAIL_VERIFICATION_FAIL"
	CodeInvalidEmailOrPhone         = "INVALID_EMAIL_OR_PHONE"
	CodeEmailAlreadyExists          = "EMAIL_ALREADY_EXISTS"
	CodePhoneAlreadyExists          = "PHONE_ALREADY_EXISTS"
	Code2FAResendFail               = "2FA_RESEND_FAIL"
	CodeVerificationUsed            = "VERIFICATION_USED"
)

const (
	CodeInvalidRequestBody    = CodeInvalidBody
	CodeAuthInvalidBody       = CodeInvalidBody
	CodeAuthInvalidEmail      = CodeInvalidEmail
	CodeAuthInvalidPhone      = CodeInvalidPhone
	CodeAuthInvalidParams     = CodeInvalidParams
	CodeInvalidPasswordLength = CodeInvalidPassword
)

// ==== Success Messages ====
const (
	MsgRegisterSuccess                 = "ثبت‌نام با موفقیت انجام شد. لطفا کد تایید ارسال شده را وارد کنید."
	MsgLoginSuccess                    = "ورود موفقیت‌آمیز بود"
	Msg2FASent                         = "کد تایید دو عاملی ارسال شد. لطفا آن را وارد کنید."
	Msg2FAVerified                     = "تایید دو عاملی با موفقیت انجام شد"
	MsgLogoutSuccess                   = "خروج با موفقیت انجام شد"
	MsgTokenRefreshed                  = "توکن جدید با موفقیت صادر شد"
	MsgForgotPasswordSent              = "در صورت وجود ایمیل یا موبایل، کد بازیابی ارسال شد"
	MsgResetPasswordSuccess            = "رمز عبور با موفقیت تغییر یافت"
	MsgResetPasswordRequestStarted     = "🔄 Reset password request"
	MsgPhoneVerificationSent           = "کد تایید موبایل ارسال شد"
	MsgPhoneVerified                   = "شماره موبایل با موفقیت تایید شد"
	MsgResendPhoneVerificationSent     = "کد تایید موبایل مجدداً ارسال شد"
	MsgEmailVerified                   = "ایمیل با موفقیت تایید شد"
	MsgEmailVerificationSent           = "کد تایید ایمیل ارسال شد"
	MsgRegisterUserCreatedBut2FAFailed = "ثبت‌نام انجام شد ولی ارسال کد تایید با خطا مواجه شد. لطفاً بعداً مجدداً تلاش کنید."
	MsgVerificationResent              = "کد تایید مجدداً ارسال شد"
	MsgInvalidOrExpiredCode            = "کد بازیابی اشتباه یا منقضی شده"
	Msg2FAResent                       = "کد تایید دو عاملی مجدداً ارسال شد"
	Msg2FAVerificationStarted          = "بررسی کد ۲FA آغاز شد"
)

// ==== Purposes for Verification Codes ====
const (
	PurposeEmailVerification = "email_verification"
	PurposePhoneVerification = "phone_verification"
	PurposeLogin2FA          = "login_2fa"
	PurposeRegister2FA       = "register_2fa"
	PurposeForgetPassword    = "forgot_password"
)

// ==== Channels for Verification Codes ====
const (
	ChannelEmail = "email"
	ChannelSMS   = "sms"
)

var ValidChannels = []string{
	ChannelEmail,
	ChannelSMS,
}

const Default2FAExpireMinutes = 2
