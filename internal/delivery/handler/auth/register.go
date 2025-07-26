package auth

import (
	"errors"
	"exchange-common/internal/consts"
	"exchange-common/internal/delivery/requestmeta"
	"exchange-common/internal/model"
	"exchange-common/internal/richerror"
	"exchange-common/internal/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Register godoc
// @Summary ثبت نام کاربر (ایمیل یا موبایل)
// @Description ایجاد حساب کاربری جدید با ایمیل یا موبایل، نام کامل، رمز عبور و کپچا
// @Tags auth
// @Accept json
// @Produce json
// @Param body body model.RegisterRequest true "اطلاعات ثبت ‌نام شامل ایمیل یا موبایل، نام کامل، رمز عبور و کپچا"
// @Success 201 {object} model.Response[model.RegisterResponse] "ثبت ‌نام موفق"
// @Failure 400 {object} model.ErrorResponseStruct "ورودی نامعتبر یا کپچا اشتباه"
// @Failure 409 {object} model.ErrorResponseStruct "کاربر قبلاً ثبت شده"
// @Failure 500 {object} model.ErrorResponseStruct "خطای داخلی سرور"
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	const op = consts.OpAuthRegister
	meta := requestmeta.NewRequestMeta(c)

	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.Handle(c, op,
			consts.ErrAuthInvalidBody,
			consts.CodeInvalidRequestBody,
			richerror.KindInvalid, err)
		return
	}

	// --- تعیین نوع شناسه: ایمیل یا موبایل ---
	var (
		email, phone string
		channel      string
	)
	switch {
	case util.ValidateEmail(req.Identifier):
		email = util.NormalizeEmail(req.Identifier)
		channel = consts.ChannelEmail
	case util.ValidatePhone(req.Identifier):
		phone = util.NormalizePhone(req.Identifier)
		channel = consts.ChannelSMS
	default:
		richerror.Handle(c, op,
			consts.ErrInvalidEmailOrPhone,
			consts.CodeInvalidEmailOrPhone,
			richerror.KindInvalid, nil)
		return
	}

	// --- بررسی رمز عبور ---
	if err := util.ValidatePassword(req.Password); err != nil {
		richerror.Handle(c, op,
			consts.ErrInvalidPassword,
			consts.CodeInvalidPasswordLength,
			richerror.KindInvalid, err)
		return
	}

	// --- بررسی کپچا ---
	if h.verifyCaptchaWrapper(c, op, c.Request.Context(), h.captchaStore, req.CaptchaID, req.CaptchaAns) {
		return
	}

	meta.Logger.Info("📝 شروع فرآیند ثبت‌نام",
		zap.String("email", email),
		zap.String("phone", phone),
	)

	// --- ثبت کاربر ---
	userID, err := h.authSvc.RegisterUser(c.Request.Context(), email, phone, req.FullName, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmailExists):
			richerror.HandleWrap(c, op,
				consts.ErrEmailExists,
				consts.CodeEmailAlreadyExists,
				richerror.KindConflict, err)
		case errors.Is(err, model.ErrPhoneExists):
			richerror.HandleWrap(c, op,
				consts.ErrPhoneExists,
				consts.CodePhoneAlreadyExists,
				richerror.KindConflict, err)
		case errors.Is(err, model.ErrUserExists):
			richerror.HandleWrap(c, op,
				consts.ErrAuthUserExists,
				consts.CodeUserAlreadyExists,
				richerror.KindConflict, err)
		default:
			richerror.HandleWrap(c, op,
				"خطا در ثبت کاربر",
				"REGISTER_FAIL",
				richerror.KindInternal, err)
		}
		return
	}

	// --- ارسال کد تأیید دو عاملی ---
	if err := h.verificationSVC.Send2FACode(c, *userID, req.Identifier, channel, consts.PurposeRegister2FA); err != nil {

		// حذف کاربر در صورت شکست در ارسال کد
		if deleteErr := h.authSvc.DeleteUser(c.Request.Context(), *userID); deleteErr != nil {
			meta.Logger.Error("❌ حذف کاربر پس از شکست ارسال کد ۲FA با خطا مواجه شد",
				zap.String("user_id", userID.String()),
				zap.Error(deleteErr),
			)
		}

		meta.Logger.Error("❌ ارسال کد ۲FA با خطا مواجه شد",
			zap.String("user_id", userID.String()),
			zap.String("channel", channel),
			zap.Error(err),
		)

		richerror.HandleWrap(c, op,
			consts.Err2FASendFail,
			consts.Code2FASendFail,
			richerror.KindInternal, err)
		return
	}

	// --- موفقیت ---
	meta.Logger.Info("✅ ثبت‌نام و ارسال کد ۲FA موفقیت‌آمیز بود",
		zap.String("user_id", userID.String()),
		zap.String("channel", channel),
		zap.Duration("duration", meta.Elapsed()),
	)

	model.SuccessResponse(c, http.StatusCreated, model.RegisterResponse{
		UserID:  userID.String(),
		Message: consts.MsgRegisterSuccess,
	})
}
