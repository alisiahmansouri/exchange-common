package auth

import (
	"exchange-common/internal/consts"
	"exchange-common/internal/delivery/requestmeta"
	"exchange-common/internal/model"
	"exchange-common/internal/richerror"
	"exchange-common/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	maxAttempt = 5
	expiry     = 15 * time.Minute
)

// ---------- Register 2FA ----------

// VerifyRegister2FA godoc
// @Summary      تایید کد دو عاملی ثبت‌نام (2FA)
// @Description  بررسی کد تایید ثبت‌نام (از طریق ایمیل یا پیامک) و صدور توکن در صورت موفقیت
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body   model.Verify2FARequest  true  "اطلاعات تایید 2FA"
// @Success      200   {object}  model.Response[model.LoginResponse]  "کد تایید صحیح بود، توکن‌ها صادر شدند"
// @Failure      400   {object}  model.ErrorResponseStruct  "بدنه درخواست نامعتبر یا کد نامعتبر"
// @Failure      401   {object}  model.ErrorResponseStruct  "کد تایید اشتباه یا منقضی شده"
// @Failure      429   {object}  model.ErrorResponseStruct  "تلاش بیش از حد مجاز"
// @Failure      500   {object}  model.ErrorResponseStruct  "خطای داخلی سرور"
// @Router       /auth/verify-register-2fa [post]
func (h *Handler) VerifyRegister2FA(c *gin.Context) {
	const (
		op      = consts.OpAuthVerifyRegister2FA
		purpose = consts.PurposeRegister2FA
	)
	meta := requestmeta.NewRequestMeta(c)

	var req model.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.Handle(c, op, consts.ErrInvalidBody, consts.CodeInvalidBody, richerror.KindInvalid, err)
		return
	}
	req.Code = strings.TrimSpace(req.Code)
	req.Identifier = strings.TrimSpace(req.Identifier)

	if req.Code == "" {
		richerror.Handle(c, op, consts.Err2FACodeInvalid, consts.CodeInvalid2FACode, richerror.KindInvalid, nil)
		return
	}
	if req.Identifier == "" {
		richerror.Handle(c, op, consts.ErrInvalidEmailOrPhone, consts.CodeInvalidEmailOrPhone, richerror.KindInvalid, nil)
		return
	}

	user, err := h.authSvc.FindUserByIdentifier(c.Request.Context(), req.Identifier)
	if err != nil {
		meta.Logger.Error("❌ "+consts.ErrFindUserFail, zap.Error(err))
		richerror.HandleWrap(c, op, consts.ErrFindUserFail, consts.CodeInvalidCredentials, richerror.KindUnauthorized, err)
		return
	}
	if user == nil {
		richerror.Handle(c, op, consts.ErrUserNotFound, consts.CodeInvalidCredentials, richerror.KindUnauthorized, nil)
		return
	}

	key := h.get2FAKey(user.ID, purpose)
	attempt, err := h.get2FAAttemptCount(c, key)
	if err != nil {
		meta.Logger.Error("❌ Redis 2FA attempt count failed", zap.Error(err))
		richerror.HandleWrap(c, op, consts.Err2FAAttemptCheckFail, consts.Code2FAAttemptCheckFail, richerror.KindInternal, err)
		return
	}
	if attempt >= maxAttempt {
		meta.Logger.Warn("🚫 "+consts.ErrTooManyAttempts, zap.String("user_id", user.ID.String()))
		richerror.Handle(c, op, consts.ErrTooManyAttempts, consts.CodeTooManyAttempts, richerror.KindTooManyRequests, nil)
		return
	}

	channel := util.Get2FAChannelByIdentifier(req.Identifier)
	if channel == "" {
		richerror.Handle(c, op, consts.ErrInvalidChannel, consts.CodeInvalidChannel, richerror.KindInvalid, nil)
		return
	}

	if err := h.verificationSVC.VerifyCode(c.Request.Context(), user.ID, channel, purpose, req.Code); err != nil {
		_ = h.increment2FAAttempt(c, key, expiry)
		meta.Logger.Warn("❌ 2FA code invalid", zap.String("user_id", user.ID.String()))
		richerror.HandleWrap(c, op, consts.Err2FACodeInvalid, consts.CodeInvalid2FACode, richerror.KindUnauthorized, err)
		return
	}
	if err := h.redisClient.Del(c.Request.Context(), key).Err(); err != nil {
		meta.Logger.Warn("❗ Redis DEL 2FA attempt key failed", zap.Error(err))
	}

	// ست کردن verified بودن ایمیل یا موبایل
	if util.ValidateEmail(req.Identifier) {
		if err := h.authSvc.MarkEmailVerified(c.Request.Context(), user.ID); err != nil {
			meta.Logger.Warn("ایمیل به عنوان verified ست نشد", zap.Error(err))
		}
	} else if util.ValidatePhone(req.Identifier) {
		if err := h.authSvc.MarkPhoneVerified(c.Request.Context(), user.ID); err != nil {
			meta.Logger.Warn("موبایل به عنوان verified ست نشد", zap.Error(err))
		}
	}

	accessToken, refreshToken, errSent := h.generateTokens(c, op, user.ID)
	if errSent {
		meta.Logger.Error("⛔ Token generation failed", zap.String("user_id", user.ID.String()))
		richerror.Handle(c, op, consts.ErrAuthTokenGenFail, consts.CodeAuthTokenGenFail, richerror.KindInternal, nil)
		return
	}

	meta.Logger.Info("✅ 2FA register verified & tokens issued", zap.String("user_id", user.ID.String()))
	model.SuccessResponse(c, http.StatusOK, model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, consts.Msg2FAVerified)
}

// ---------- Login 2FA ----------

// VerifyLogin2FA godoc
// @Summary      تایید کد دو عاملی ورود (2FA)
// @Description  بررسی کد تایید ارسال‌شده هنگام ورود و صدور توکن جدید در صورت موفقیت
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body   model.Verify2FARequest  true  "کد تایید و شناسه کاربر"
// @Success      200   {object}  model.Response[model.LoginResponse]  "کد تایید صحیح بود، توکن صادر شد"
// @Failure      400   {object}  model.ErrorResponseStruct  "بدنه درخواست نامعتبر یا کد خالی"
// @Failure      401   {object}  model.ErrorResponseStruct  "کد تایید اشتباه یا منقضی شده"
// @Failure      429   {object}  model.ErrorResponseStruct  "تعداد تلاش‌های بیش از حد مجاز"
// @Failure      500   {object}  model.ErrorResponseStruct  "خطای داخلی سرور"
// @Router       /auth/verify-login-2fa [post]
func (h *Handler) VerifyLogin2FA(c *gin.Context) {
	const (
		op      = consts.OpAuthVerifyLogin2FA
		purpose = consts.PurposeLogin2FA
	)

	meta := requestmeta.NewRequestMeta(c)
	var req model.Verify2FARequest

	// 1. اعتبارسنجی بدنه ورودی
	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.Handle(c, op, consts.ErrInvalidBody, consts.CodeInvalidBody, richerror.KindInvalid, err)
		return
	}
	req.Code = strings.TrimSpace(req.Code)
	req.Identifier = strings.TrimSpace(req.Identifier)

	if req.Code == "" || req.Identifier == "" {
		richerror.Handle(c, op, consts.Err2FACodeInvalid, consts.CodeInvalid2FACode, richerror.KindInvalid, nil)
		return
	}

	// 2. یافتن کاربر
	user, err := h.authSvc.FindUserByIdentifier(c.Request.Context(), req.Identifier)
	if err != nil {
		meta.Logger.Error("❌ "+consts.ErrFindUserFail, zap.Error(err))
		richerror.HandleWrap(c, op, consts.ErrFindUserFail, consts.CodeInvalidCredentials, richerror.KindUnauthorized, err)
		return
	}
	if user == nil {
		richerror.Handle(c, op, consts.ErrUserNotFound, consts.CodeInvalidCredentials, richerror.KindUnauthorized, nil)
		return
	}

	// 3. محدودیت تعداد تلاش
	key := h.get2FAKey(user.ID, purpose)
	attempt, err := h.get2FAAttemptCount(c, key)
	if err != nil {
		meta.Logger.Error("❌ Redis 2FA attempt count failed", zap.Error(err))
		richerror.HandleWrap(c, op, consts.Err2FAAttemptCheckFail, consts.Code2FAAttemptCheckFail, richerror.KindInternal, err)
		return
	}
	if attempt >= maxAttempt {
		meta.Logger.Warn("🚫 "+consts.ErrTooManyAttempts, zap.String("user_id", user.ID.String()))
		richerror.Handle(c, op, consts.ErrTooManyAttempts, consts.CodeTooManyAttempts, richerror.KindTooManyRequests, nil)
		return
	}

	// 4. تایید کد ۲FA
	channel := util.Get2FAChannelByIdentifier(req.Identifier)
	if channel == "" {
		richerror.Handle(c, op, consts.ErrInvalidChannel, consts.CodeInvalidChannel, richerror.KindInvalid, nil)
		return
	}
	if err := h.verificationSVC.VerifyCode(
		c.Request.Context(),
		user.ID,
		channel,
		purpose,
		req.Code,
	); err != nil {
		_ = h.increment2FAAttempt(c, key, expiry)
		meta.Logger.Warn("❌ 2FA code invalid", zap.String("user_id", user.ID.String()))
		richerror.HandleWrap(c, op, consts.Err2FACodeInvalid, consts.CodeInvalid2FACode, richerror.KindUnauthorized, err)
		return
	}
	if err := h.redisClient.Del(c.Request.Context(), key).Err(); err != nil {
		meta.Logger.Warn("❗ Redis DEL 2FA attempt key failed", zap.Error(err))
	}

	// 5. صدور توکن و ریترن
	accessToken, refreshToken, errSent := h.generateTokens(c, op, user.ID)
	if errSent {
		meta.Logger.Error("⛔ Token generation failed", zap.String("user_id", user.ID.String()))
		richerror.Handle(c, op, consts.ErrAuthTokenGenFail, consts.CodeAuthTokenGenFail, richerror.KindInternal, nil)
		return
	}

	meta.Logger.Info("✅ 2FA login verified & tokens issued", zap.String("user_id", user.ID.String()))
	model.SuccessResponse(c, http.StatusOK, model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, consts.Msg2FAVerified)
}

// ---------- Email Verify ----------

// VerifyEmail2FA godoc
// @Summary      تایید ایمیل با کد ارسالی
// @Description  بررسی کد تایید ایمیل ارسال‌شده به کاربر و ثبت تایید ایمیل (مناسب تایید بعد از ثبت‌نام یا درخواست کاربر)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body   model.Verify2FARequest  true  "کد تایید و شناسه (ایمیل)"
// @Success      200   {object}  model.Response[model.SimpleMessageResponse]  "ایمیل با موفقیت تایید شد"
// @Failure      400   {object}  model.ErrorResponseStruct  "بدنه درخواست نامعتبر یا کد خالی"
// @Failure      401   {object}  model.ErrorResponseStruct  "کد تایید اشتباه یا منقضی شده"
// @Failure      429   {object}  model.ErrorResponseStruct  "تعداد تلاش‌های بیش از حد مجاز"
// @Failure      500   {object}  model.ErrorResponseStruct  "خطای داخلی سرور"
// @Router       /auth/verify-email [post]
func (h *Handler) VerifyEmail2FA(c *gin.Context) {
	const (
		op         = consts.OpAuthVerifyEmail
		purpose    = consts.PurposeEmailVerification
		maxAttempt = 5
		expiry     = 15 * time.Minute
	)
	meta := requestmeta.NewRequestMeta(c)
	var req model.Verify2FARequest

	// 1. اعتبارسنجی بدنه
	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.Handle(c, op, consts.ErrInvalidBody, consts.CodeInvalidBody, richerror.KindInvalid, err)
		return
	}
	req.Code = strings.TrimSpace(req.Code)
	req.Identifier = strings.TrimSpace(req.Identifier)
	if req.Code == "" || req.Identifier == "" {
		richerror.Handle(c, op, consts.Err2FACodeInvalid, consts.CodeInvalid2FACode, richerror.KindInvalid, nil)
		return
	}

	// 2. فقط ایمیل را قبول کن
	if !util.ValidateEmail(req.Identifier) {
		richerror.Handle(c, op, consts.ErrInvalidEmail, consts.CodeInvalidEmail, richerror.KindInvalid, nil)
		return
	}

	// 3. یافتن کاربر
	user, err := h.authSvc.FindUserByIdentifier(c.Request.Context(), req.Identifier)
	if err != nil {
		meta.Logger.Error("❌ "+consts.ErrFindUserFail, zap.Error(err))
		richerror.HandleWrap(c, op, consts.ErrFindUserFail, consts.CodeInvalidCredentials, richerror.KindUnauthorized, err)
		return
	}
	if user == nil {
		richerror.Handle(c, op, consts.ErrUserNotFound, consts.CodeInvalidCredentials, richerror.KindUnauthorized, nil)
		return
	}

	// 4. محدودیت تلاش
	key := h.get2FAKey(user.ID, purpose)
	attempt, err := h.get2FAAttemptCount(c, key)
	if err != nil {
		meta.Logger.Error("❌ Redis 2FA attempt count failed", zap.Error(err))
		richerror.HandleWrap(c, op, consts.Err2FAAttemptCheckFail, consts.Code2FAAttemptCheckFail, richerror.KindInternal, err)
		return
	}
	if attempt >= maxAttempt {
		meta.Logger.Warn("🚫 "+consts.ErrTooManyAttempts, zap.String("user_id", user.ID.String()))
		richerror.Handle(c, op, consts.ErrTooManyAttempts, consts.CodeTooManyAttempts, richerror.KindTooManyRequests, nil)
		return
	}

	// 5. بررسی صحت کد
	channel := util.Get2FAChannelByIdentifier(req.Identifier)
	if channel != consts.ChannelEmail {
		richerror.Handle(c, op, consts.ErrInvalidChannel, consts.CodeInvalidChannel, richerror.KindInvalid, nil)
		return
	}
	if err := h.verificationSVC.VerifyCode(
		c.Request.Context(),
		user.ID,
		channel,
		purpose,
		req.Code,
	); err != nil {
		_ = h.increment2FAAttempt(c, key, expiry)
		meta.Logger.Warn("❌ 2FA code invalid", zap.String("user_id", user.ID.String()))
		richerror.HandleWrap(c, op, consts.Err2FACodeInvalid, consts.CodeInvalid2FACode, richerror.KindUnauthorized, err)
		return
	}
	if err := h.redisClient.Del(c.Request.Context(), key).Err(); err != nil {
		meta.Logger.Warn("❗ Redis DEL 2FA attempt key failed", zap.Error(err))
	}

	// 6. ست کردن verified
	if err := h.authSvc.MarkEmailVerified(c.Request.Context(), user.ID); err != nil {
		meta.Logger.Warn("ایمیل به عنوان verified ست نشد", zap.Error(err))
	}

	meta.Logger.Info("✅ Email verified", zap.String("user_id", user.ID.String()))

	model.SuccessResponse(c, http.StatusOK, model.SimpleMessageResponse{
		Success: true,
		Message: consts.MsgEmailVerified,
	})
}

// ---------- Phone Verify ----------

// VerifyPhone2FA godoc
// @Summary      تایید شماره موبایل با کد ارسالی
// @Description  بررسی کد تایید ارسال‌شده به شماره موبایل و ثبت تایید شماره موبایل (پس از ثبت‌نام یا ویرایش شماره)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body   model.Verify2FARequest  true  "کد تایید و شناسه (موبایل)"
// @Success      200   {object}  model.Response[model.SimpleMessageResponse]  "شماره موبایل با موفقیت تایید شد"
// @Failure      400   {object}  model.ErrorResponseStruct  "بدنه درخواست یا کد تایید نامعتبر"
// @Failure      401   {object}  model.ErrorResponseStruct  "کد تایید اشتباه یا منقضی شده"
// @Failure      429   {object}  model.ErrorResponseStruct  "تعداد تلاش‌های بیش از حد مجاز"
// @Failure      500   {object}  model.ErrorResponseStruct  "خطای داخلی سرور"
// @Router       /auth/verify-phone [post]
func (h *Handler) VerifyPhone2FA(c *gin.Context) {
	const (
		op         = consts.OpAuthVerifyPhone
		purpose    = consts.PurposePhoneVerification
		maxAttempt = 5
		expiry     = 15 * time.Minute
	)
	meta := requestmeta.NewRequestMeta(c)
	var req model.Verify2FARequest

	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.Handle(c, op, consts.ErrInvalidBody, consts.CodeInvalidBody, richerror.KindInvalid, err)
		return
	}
	req.Code = strings.TrimSpace(req.Code)
	req.Identifier = strings.TrimSpace(req.Identifier)
	if req.Code == "" || req.Identifier == "" {
		richerror.Handle(c, op, consts.Err2FACodeInvalid, consts.CodeInvalid2FACode, richerror.KindInvalid, nil)
		return
	}
	if !util.ValidatePhone(req.Identifier) {
		richerror.Handle(c, op, consts.ErrAuthInvalidPhone, consts.CodeInvalidPhone, richerror.KindInvalid, nil)
		return
	}

	user, err := h.authSvc.FindUserByIdentifier(c.Request.Context(), req.Identifier)
	if err != nil {
		meta.Logger.Error("❌ "+consts.ErrFindUserFail, zap.Error(err))
		richerror.HandleWrap(c, op, consts.ErrFindUserFail, consts.CodeInvalidCredentials, richerror.KindUnauthorized, err)
		return
	}
	if user == nil {
		richerror.Handle(c, op, consts.ErrUserNotFound, consts.CodeInvalidCredentials, richerror.KindUnauthorized, nil)
		return
	}

	key := h.get2FAKey(user.ID, purpose)
	attempt, err := h.get2FAAttemptCount(c, key)
	if err != nil {
		meta.Logger.Error("❌ Redis 2FA attempt count failed", zap.Error(err))
		richerror.HandleWrap(c, op, consts.Err2FAAttemptCheckFail, consts.Code2FAAttemptCheckFail, richerror.KindInternal, err)
		return
	}
	if attempt >= maxAttempt {
		meta.Logger.Warn("🚫 "+consts.ErrTooManyAttempts, zap.String("user_id", user.ID.String()))
		richerror.Handle(c, op, consts.ErrTooManyAttempts, consts.CodeTooManyAttempts, richerror.KindTooManyRequests, nil)
		return
	}

	channel := util.Get2FAChannelByIdentifier(req.Identifier)
	if channel != consts.ChannelSMS {
		richerror.Handle(c, op, consts.ErrInvalidChannel, consts.CodeInvalidChannel, richerror.KindInvalid, nil)
		return
	}
	if err := h.verificationSVC.VerifyCode(
		c.Request.Context(),
		user.ID,
		channel,
		purpose,
		req.Code,
	); err != nil {
		_ = h.increment2FAAttempt(c, key, expiry)
		meta.Logger.Warn("❌ 2FA code invalid", zap.String("user_id", user.ID.String()))
		richerror.HandleWrap(c, op, consts.Err2FACodeInvalid, consts.CodeInvalid2FACode, richerror.KindUnauthorized, err)
		return
	}
	if err := h.redisClient.Del(c.Request.Context(), key).Err(); err != nil {
		meta.Logger.Warn("❗ Redis DEL 2FA attempt key failed", zap.Error(err))
	}

	if err := h.authSvc.MarkPhoneVerified(c.Request.Context(), user.ID); err != nil {
		meta.Logger.Warn("شماره موبایل به عنوان verified ست نشد", zap.Error(err))
	}

	meta.Logger.Info("✅ Phone verified", zap.String("user_id", user.ID.String()))

	model.SuccessResponse(c, http.StatusOK, model.SimpleMessageResponse{
		Success: true,
		Message: consts.MsgPhoneVerified,
	})
}

// ------ helpers -------
func (h *Handler) get2FAAttemptCount(c *gin.Context, key string) (int, error) {
	val, err := h.redisClient.Get(c.Request.Context(), key).Result()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}

func (h *Handler) increment2FAAttempt(c *gin.Context, key string, expire time.Duration) error {
	pipe := h.redisClient.TxPipeline()
	pipe.Incr(c.Request.Context(), key)
	pipe.Expire(c.Request.Context(), key, expire)
	_, err := pipe.Exec(c.Request.Context())
	return err
}

func (h *Handler) get2FAKey(userID uuid.UUID, purpose string) string {
	return fmt.Sprintf("2fa_attempts:%s:%s", userID.String(), purpose)
}

// ResendVerification godoc
// @Summary ارسال مجدد کد تایید ایمیل یا موبایل
// @Description ارسال دوباره کد تایید به شماره موبایل یا ایمیل کاربر (با محدودیت نرخ)
// @Tags auth
// @Accept json
// @Produce json
// @Param body body model.ResendVerificationRequest true "ایمیل یا شماره موبایل"
// @Success 200 {object} model.Response[model.SimpleMessageResponse] "کد تایید مجدداً ارسال شد"
// @Failure 400 {object} model.ErrorResponseStruct "ورودی نامعتبر"
// @Failure 404 {object} model.ErrorResponseStruct "کاربر یافت نشد"
// @Failure 429 {object} model.ErrorResponseStruct "محدودیت ارسال کد بیش از حد مجاز است"
// @Failure 500 {object} model.ErrorResponseStruct "خطای داخلی سرور"
// @Router /auth/resend-verification [post]
func (h *Handler) ResendVerification(c *gin.Context) {
	const op = consts.OpAuthResendVerification
	meta := requestmeta.NewRequestMeta(c)

	var req model.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.Handle(c, op, consts.ErrInvalidBody, consts.CodeInvalidBody, richerror.KindInvalid, err)
		return
	}

	identifier := strings.TrimSpace(req.Identifier)
	ctx := c.Request.Context()

	if identifier == "" {
		richerror.Handle(c, op, consts.ErrInvalidEmailOrPhone, consts.CodeInvalidEmailOrPhone, richerror.KindInvalid, nil)
		return
	}

	canSend, err := h.canSendVerification(ctx, identifier)
	if err != nil {
		meta.Logger.Error("❌ "+consts.ErrRateLimitCheckFail, zap.Error(err))
		richerror.HandleWrap(c, op, consts.ErrRateLimitCheckFail, consts.CodeRateLimitCheckFail, richerror.KindInternal, err)
		return
	}
	if !canSend {
		meta.Logger.Warn("📛 "+consts.ErrRateLimitExceeded, zap.String("identifier", identifier))
		richerror.Handle(c, op, consts.ErrRateLimitExceeded, consts.CodeRateLimitExceeded, richerror.KindTooManyRequests, nil)
		return
	}

	channel := util.Get2FAChannelByIdentifier(identifier)
	if channel == "" {
		richerror.Handle(c, op, consts.ErrInvalidEmailOrPhone, consts.CodeInvalidEmailOrPhone, richerror.KindInvalid, nil)
		return
	}

	user, err := h.authSvc.FindUserByIdentifier(ctx, identifier)
	if err != nil {
		meta.Logger.Error("❌ "+consts.ErrFindUserFail, zap.Error(err))
		richerror.HandleWrap(c, op, consts.ErrFindUserFail, consts.CodeInvalidCredentials, richerror.KindNotFound, err)
		return
	}
	if user == nil {
		richerror.Handle(c, op, consts.ErrUserNotFound, consts.CodeInvalidCredentials, richerror.KindNotFound, nil)
		return
	}

	purpose := strings.TrimSpace(req.Purpose)
	if purpose == "" {
		purpose = consts.PurposeRegister2FA
	}
	switch purpose {
	case consts.PurposeRegister2FA, consts.PurposeLogin2FA, consts.PurposeEmailVerification, consts.PurposePhoneVerification:
		// ok
	default:
		richerror.Handle(c, op, consts.ErrInvalidPurpose, consts.CodeInvalidPurpose, richerror.KindInvalid, nil)
		return
	}

	if err := h.verificationSVC.Send2FACode(ctx, user.ID, identifier, channel, purpose); err != nil {
		meta.Logger.Error("❌ "+consts.Err2FASendFail, zap.String("user_id", user.ID.String()), zap.Error(err))
		richerror.HandleWrap(c, op, consts.Err2FASendFail, consts.Code2FASendFail, richerror.KindInternal, err)
		return
	}

	meta.Logger.Info("🔁 "+consts.Msg2FAResent,
		zap.String("user_id", user.ID.String()),
		zap.String("identifier", identifier),
		zap.String("channel", channel),
		zap.String("purpose", purpose),
		zap.Duration("duration", meta.Elapsed()),
	)

	model.SuccessResponse(c, http.StatusOK, model.SimpleMessageResponse{
		Success: true,
		Message: consts.Msg2FAResent,
	})
}
