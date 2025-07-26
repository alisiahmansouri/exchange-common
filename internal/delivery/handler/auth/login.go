package auth

import (
	"exchange-common/internal/consts"
	"exchange-common/internal/delivery/requestmeta"
	"exchange-common/internal/model"
	"exchange-common/internal/richerror"
	"exchange-common/internal/util"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strings"

	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Login godoc
// @Summary ورود کاربر
// @Description احراز هویت کاربر با ایمیل یا موبایل، رمز عبور و کپچا. در صورت فعال بودن 2FA، ارسال کد تایید دو عاملی
// @Tags auth
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "اطلاعات ورود همراه کپچا"
// @Success 200 {object} model.Response[model.LoginResponse] "ورود موفق بدون 2FA"
// @Success 202 {object} model.Response[model.SimpleMessageResponse] "کد 2FA ارسال شد، منتظر تایید باشید"
// @Failure 400 {object} model.ErrorResponseStruct "بدنه درخواست یا کپچا نامعتبر"
// @Failure 401 {object} model.ErrorResponseStruct "اطلاعات ورود یا کپچا اشتباه"
// @Failure 403 {object} model.ErrorResponseStruct "ایمیل یا موبایل تایید نشده"
// @Failure 429 {object} model.ErrorResponseStruct "تعداد تلاش‌ها بیش از حد مجاز"
// @Failure 500 {object} model.ErrorResponseStruct "خطای داخلی سرور"
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	const op = consts.OpAuthLogin
	meta := requestmeta.NewRequestMeta(c)
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.Handle(c, op, consts.ErrInvalidBody, consts.CodeInvalidBody, richerror.KindInvalid, err)
		return
	}

	req.Identifier = strings.TrimSpace(req.Identifier)
	req.Password = strings.TrimSpace(req.Password)

	// 1. محدودیت تلاش لاگین
	if h.isLoginThrottled(c.Request.Context(), req.Identifier) {
		meta.Logger.Warn("⛔ Too many login attempts",
			zap.String("identifier", req.Identifier),
			zap.String("ip", c.ClientIP()),
		)
		richerror.Handle(c, op, consts.ErrTooManyAttempts, consts.CodeLoginThrottled, richerror.KindTooManyRequests, nil)
		return
	}

	// 2. اعتبارسنجی identifier
	if !util.ValidateEmail(req.Identifier) && !util.ValidatePhone(req.Identifier) {
		h.incrementLoginAttempt(c.Request.Context(), req.Identifier)
		richerror.Handle(c, op, consts.ErrInvalidEmailOrPhone, consts.CodeInvalidEmailOrPhone, richerror.KindInvalid, nil)
		return
	}

	// 3. اعتبارسنجی رمز عبور
	if err := util.ValidatePassword(req.Password); err != nil {
		h.incrementLoginAttempt(c.Request.Context(), req.Identifier)
		richerror.Handle(c, op, consts.ErrInvalidPassword, consts.CodeInvalidPassword, richerror.KindInvalid, nil)
		return
	}

	// 4. کپچا
	if h.verifyCaptchaWrapper(c, op, c.Request.Context(), h.captchaStore, req.CaptchaID, req.CaptchaAns) {
		return // کپچا اشتباه → تلاش ناموفق ثبت نشود
	}

	meta.Logger.Info("🔐 Login attempt",
		zap.String("identifier", req.Identifier),
		zap.String("ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
	)

	// 5. تلاش برای لاگین کاربر
	user, err := h.authSvc.AuthenticateUser(c.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		h.incrementLoginAttempt(c.Request.Context(), req.Identifier)
		richerror.HandleWrap(c, op, consts.ErrAuthInvalidCredentials, consts.CodeInvalidCredentials, richerror.KindUnauthorized, err)
		return
	}

	// 6. ریست تلاش‌ها
	h.resetLoginAttempts(c.Request.Context(), req.Identifier)

	// 7. چک تایید ایمیل/موبایل
	if util.ValidateEmail(req.Identifier) && !user.IsEmailVerified {
		richerror.Handle(c, op, "ایمیل شما هنوز تایید نشده است", "EMAIL_NOT_VERIFIED", richerror.KindForbidden, nil)
		return
	}
	if util.ValidatePhone(req.Identifier) && !user.IsPhoneVerified {
		richerror.Handle(c, op, "شماره موبایل شما هنوز تایید نشده است", "PHONE_NOT_VERIFIED", richerror.KindForbidden, nil)
		return
	}

	// 8. چک فعال بودن 2FA
	twoFAEnabled, err := h.authSvc.Is2FAEnabled(c.Request.Context(), user.ID)
	if err != nil {
		richerror.HandleWrap(c, op, consts.Err2FACodeCheckFail, consts.Code2FACheckError, richerror.KindInternal, err)
		return
	}

	if twoFAEnabled {
		channel := util.Get2FAChannelByIdentifier(req.Identifier)
		err := h.verificationSVC.Send2FACode(c.Request.Context(), user.ID, req.Identifier, channel, consts.PurposeLogin2FA)
		if err != nil {
			meta.Logger.Error("❌ Failed to send 2FA code", zap.String("user_id", user.ID.String()), zap.Error(err))
			richerror.HandleWrap(c, op, consts.Err2FASendFail, consts.Code2FASendFail, richerror.KindInternal, err)
			return
		}
		meta.Logger.Info("📨 2FA code sent", zap.String("user_id", user.ID.String()))
		model.SuccessResponse(c, http.StatusAccepted, model.SimpleMessageResponse{
			Success: true,
			Message: consts.Msg2FASent,
		})
		return
	}

	// 9. صدور توکن (اگر 2FA نبود)
	accessToken, refreshToken, errSent := h.generateTokens(c, op, user.ID)
	if errSent {
		meta.Logger.Error("⛔ Token generation failed", zap.String("user_id", user.ID.String()))
		richerror.Handle(c, op, consts.ErrAuthTokenGenFail, consts.CodeAuthTokenGenFail, richerror.KindInternal, nil)
		return
	}

	meta.Logger.Info("✅ Login successful",
		zap.String("user_id", user.ID.String()),
		zap.Duration("duration", meta.Elapsed()),
	)

	model.SuccessResponse(c, http.StatusOK, model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, consts.MsgLoginSuccess)
}

// helper methods:

func loginAttemptKey(identifier string) string {
	return "login_attempts:" + strings.ToLower(identifier)
}

func (h *Handler) isLoginThrottled(ctx context.Context, identifier string) bool {
	attempts, err := h.redisClient.Get(ctx, loginAttemptKey(identifier)).Int()
	if err != nil && err != redis.Nil {
		return false // skip throttling on Redis error
	}
	return attempts >= consts.MaxLoginAttempts
}

func (h *Handler) incrementLoginAttempt(ctx context.Context, identifier string) {
	key := loginAttemptKey(identifier)
	pipe := h.redisClient.TxPipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, consts.LoginAttemptDuration)
	_, _ = pipe.Exec(ctx)
}

func (h *Handler) resetLoginAttempts(ctx context.Context, identifier string) {
	_ = h.redisClient.Del(ctx, loginAttemptKey(identifier)).Err()
}
