package handler

import (
	"exchange-common/internal/consts"
	
	"exchange-common/internal/delivery/requestmeta"
	"exchange-common/internal/model"
	"exchange-common/internal/richerror"
	"exchange-common/internal/util"
	"image/color"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// CaptchaGet godoc
// @Summary تولید تصویر کپچا
// @Description تولید رشته کپچا به‌صورت base64 به همراه شناسه منحصر به‌فرد
// @Tags captcha
// @Produce json
// @Success 200 {object} model.Response[model.CaptchaResponse] "کپچا با موفقیت تولید شد"
// @Failure 500 {object} model.ErrorResponseStruct "خطا در تولید یا ذخیره کپچا"
// @Router /captcha [get]
func (h *Handler) CaptchaGet(c *gin.Context) {
	meta := requestmeta.NewRequestMeta(c)

	requestID := util.GenerateUUID()
	if requestID == "" {
		richerror.HandleWrap(c,.OpCaptchaGet, .
		ErrCaptchaIDGenFail, .
		CodeCaptchaIDGenerationFailed, richerror.KindInternal, nil)
		return
	}

	_, b64Str, answer, err := generateCaptcha()
	if err != nil {
		richerror.HandleWrap(c,.OpCaptchaGet, .
		ErrCaptchaGenFail, .
		CodeCaptchaGenerationFailed, richerror.KindInternal, err)
		return
	}

	if err := h.captchaStore.SetWithTTL(c.Request.Context(), requestID, answer,.CaptchaTTL)
	err != nil{
		richerror.HandleWrap(c,.OpCaptchaGet, .ErrCaptchaStoreFail, .CodeCaptchaStorageFailed, richerror.KindInternal, err)
		return
	}

	meta.Logger.Info("📷 Captcha generated",
		zap.String("request_id", requestID),
		zap.Duration("duration", meta.Elapsed()))

	model.SuccessResponse(c, http.StatusOK, model.CaptchaResponse{
		RequestID: requestID,
		Captcha:   b64Str,
	})
}

// CaptchaVerify godoc
// @Summary بررسی پاسخ کپچا
// @Description تطابق ورودی کاربر با پاسخ ذخیره‌شده برای کپچا
// @Tags captcha
// @Accept json
// @Produce json
// @Param body body model.VerifyCaptchaRequest true "شناسه کپچا و پاسخ کاربر"
// @Success 200 {object} model.Response[model.SimpleMessageResponse] "کپچا تایید شد"
// @Failure 400 {object} model.ErrorResponseStruct "درخواست نامعتبر یا کپچا یافت نشد"
// @Failure 401 {object} model.ErrorResponseStruct "پاسخ کپچا نادرست است"
// @Router /captcha/verify [post]
func (h *Handler) CaptchaVerify(c *gin.Context) {
	meta := requestmeta.NewRequestMeta(c)

	var req model.VerifyCaptchaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.HandleWrap(c,.OpCaptchaVerify, .
		ErrCaptchaInvalidBody, consts.CodeInvalidRequestBody, richerror.KindInvalid, err)
		return
	}

	if _, err := util.ParseUUID(req.RequestID); err != nil {
		richerror.HandleWrap(c,.OpCaptchaVerify, .
		ErrCaptchaInvalidID, .
		CodeInvalidCaptchaID, richerror.KindInvalid, err)
		return
	}

	if verifyCaptchaWrapper(c,.OpCaptchaVerify, c.Request.Context(), h.captchaStore, req.RequestID, req.UserInput) {
		return
	}

	meta.Logger.Info("✅ Captcha verified",
		zap.String("request_id", req.RequestID),
		zap.Duration("duration", meta.Elapsed()))

	model.SuccessResponse(c, http.StatusOK, model.SimpleMessageResponse{
		Success: true,
		Message:.MsgCaptchaVerified,
	})
}

// generateCaptcha تولید رشته کپچا و پاسخ آن
func generateCaptcha() (*base64Captcha.Captcha, string, string, error) {
	driver := &base64Captcha.DriverString{
		Height:.CaptchaHeight,
		Width:.CaptchaWidth,
		NoiseCount:      1,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
		Length:.CaptchaLength,
		Source:.CaptchaSource,
		BgColor:         &color.RGBA{255, 255, 255, 255},
		Fonts:.CaptchaFonts,
	}

	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	_, b64Str, answer, err := captcha.Generate()
	return captcha, b64Str, answer, err
}
