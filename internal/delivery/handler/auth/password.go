package auth

import (
	"exchange-common/internal/consts"
	"exchange-common/internal/delivery/requestmeta"
	"exchange-common/internal/model"
	"exchange-common/internal/richerror"
	"exchange-common/internal/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// ForgotPassword godoc
// @Summary      فراموشی رمز عبور
// @Description  ارسال لینک یا کد بازیابی رمز عبور به ایمیل یا موبایل کاربر
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body model.ForgotPasswordRequest true "ایمیل یا شماره موبایل کاربر"
// @Success      200 {object} model.Response[model.SimpleMessageResponse] "در صورت موفقیت، پیام ارسال می‌شود"
// @Failure      400 {object} model.ErrorResponseStruct "ورودی نامعتبر"
// @Failure      404 {object} model.ErrorResponseStruct "ایمیل یا موبایل پیدا نشد"
// @Failure      500 {object} model.ErrorResponseStruct "خطای داخلی سرور"
// @Router       /auth/forgot-password [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	meta := requestmeta.NewRequestMeta(c)
	var req model.ForgotPasswordRequest

	// 1. Parse & validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.Handle(c, consts.OpAuthForgotPassword, consts.ErrInvalidBody, consts.CodeInvalidBody, richerror.KindInvalid, err)
		return
	}

	identifier := strings.TrimSpace(req.Identifier)
	if identifier == "" {
		richerror.Handle(c, consts.OpAuthForgotPassword, consts.ErrInvalidEmailOrPhone, consts.CodeInvalidEmailOrPhone, richerror.KindInvalid, nil)
		return
	}
	if !util.ValidateEmail(identifier) && !util.ValidatePhone(identifier) {
		richerror.Handle(c, consts.OpAuthForgotPassword, consts.ErrInvalidEmailOrPhone, consts.CodeInvalidEmailOrPhone, richerror.KindInvalid, nil)
		return
	}

	meta.Logger.Info("🔄 Forgot password request",
		zap.String("identifier", identifier),
	)

	// 2. پیدا کردن کاربر (بر اساس ایمیل یا موبایل)
	user, err := h.authSvc.FindUserByIdentifier(c.Request.Context(), identifier)
	if err != nil {
		richerror.HandleWrap(c, consts.OpAuthForgotPassword, consts.ErrFindUserFail, consts.CodeForgotPasswordEmailNotFound, richerror.KindNotFound, err)
		return
	}
	if user == nil {
		richerror.Handle(c, consts.OpAuthForgotPassword, consts.ErrForgotPasswordEmailNotFound, consts.CodeForgotPasswordEmailNotFound, richerror.KindNotFound, nil)
		return
	}

	// 4. ارسال کد یا لینک بازیابی (به ایمیل یا موبایل)
	channel := util.Get2FAChannelByIdentifier(identifier)
	purpose := consts.PurposeForgetPassword
	if err := h.verificationSVC.Send2FACode(c.Request.Context(), user.ID, req.Identifier, channel, purpose); err != nil {
		richerror.HandleWrap(c, consts.OpAuthForgotPassword, consts.ErrForgotPasswordSendFail, consts.CodeForgotPasswordSendFail, richerror.KindInternal, err)
		return
	}

	model.SuccessResponse(c, http.StatusOK, model.SimpleMessageResponse{
		Success: true,
		Message: consts.MsgForgotPasswordSent,
	})
}

// ResetPassword godoc
// @Summary      ریست رمز عبور
// @Description  بازنشانی رمز عبور با کد ارسالی به ایمیل یا موبایل
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body model.ResetPasswordRequest true "ایمیل یا موبایل، کد و رمز جدید"
// @Success      200 {object} model.Response[model.SimpleMessageResponse] "ریست موفق"
// @Failure      400 {object} model.ErrorResponseStruct "ورودی نامعتبر"
// @Failure      401 {object} model.ErrorResponseStruct "کد نامعتبر یا منقضی شده"
// @Failure      500 {object} model.ErrorResponseStruct "خطای داخلی سرور"
// @Router       /auth/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	meta := requestmeta.NewRequestMeta(c)
	var req model.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.Handle(
			c, consts.OpAuthResetPassword,
			consts.ErrInvalidBody, consts.CodeInvalidBody,
			richerror.KindInvalid, err,
		)
		return
	}

	identifier := strings.TrimSpace(req.Identifier)
	if identifier == "" || (!util.ValidateEmail(identifier) && !util.ValidatePhone(identifier)) {
		richerror.Handle(
			c, consts.OpAuthResetPassword,
			consts.ErrInvalidEmailOrPhone, consts.CodeInvalidEmailOrPhone,
			richerror.KindInvalid, nil,
		)
		return
	}

	// ۱. پیدا کردن کاربر
	user, err := h.authSvc.FindUserByIdentifier(c.Request.Context(), identifier)
	if err != nil {
		richerror.HandleWrap(
			c, consts.OpAuthResetPassword,
			consts.ErrFindUserFail, consts.CodeForgotPasswordEmailNotFound,
			richerror.KindNotFound, err,
		)
		return
	}
	if user == nil {
		richerror.Handle(
			c, consts.OpAuthResetPassword,
			consts.ErrForgotPasswordEmailNotFound, consts.CodeForgotPasswordEmailNotFound,
			richerror.KindNotFound, nil,
		)
		return
	}

	// ۲. اعتبارسنجی رمز جدید
	if err := util.ValidatePassword(req.NewPassword); err != nil {
		richerror.Handle(
			c, consts.OpAuthResetPassword,
			consts.ErrInvalidPassword, consts.CodeInvalidPassword,
			richerror.KindInvalid, nil,
		)
		return
	}

	// ۳. صحت‌سنجی کد بازیابی با توجه به کانال و پرپس
	channel := util.Get2FAChannelByIdentifier(identifier)
	purpose := consts.PurposeForgetPassword

	if err := h.verificationSVC.VerifyCode(
		c.Request.Context(),
		user.ID,
		channel,
		purpose,
		req.Code,
	); err != nil {
		richerror.HandleWrap(
			c, consts.OpAuthResetPassword,
			consts.ErrInvalidOrExpiredCode, consts.CodeInvalidOrExpiredCode,
			richerror.KindUnauthorized, err,
		)
		return
	}

	// ۴. ریست رمز عبور (بدون چک مجدد کد)
	if err := h.authSvc.ResetPassword(c.Request.Context(), user.ID, req.NewPassword); err != nil {
		richerror.HandleWrap(
			c, consts.OpAuthResetPassword,
			consts.ErrResetPasswordFail, consts.CodeResetPasswordFail,
			richerror.KindInternal, err,
		)
		return
	}

	meta.Logger.Info(consts.MsgResetPasswordSuccess, zap.String("user_id", user.ID.String()))

	model.SuccessResponse(c, http.StatusOK, model.SimpleMessageResponse{
		Success: true,
		Message: consts.MsgResetPasswordSuccess,
	})
}
