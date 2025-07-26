package auth

import (
	"exchange-common/internal/consts"

	"exchange-common/internal/delivery/requestmeta"
	"exchange-common/internal/model"
	"exchange-common/internal/richerror"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

// RefreshToken godoc
// @Summary نوسازی توکن دسترسی
// @Description صدور توکن دسترسی جدید با استفاده از توکن refresh معتبر
// @Tags auth
// @Accept json
// @Produce json
// @Param body body model.RefreshTokenRequest true "اطلاعات مربوط به توکن refresh"
// @Success 200 {object} model.Response[model.TokenResponse] "توکن با موفقیت نوسازی شد"
// @Failure 400 {object} model.ErrorResponseStruct "بدنه درخواست نامعتبر"
// @Failure 401 {object} model.ErrorResponseStruct "توکن نامعتبر یا منقضی شده"
// @Failure 500 {object} model.ErrorResponseStruct "خطای داخلی سرور"
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	meta := requestmeta.NewRequestMeta(c)
	var req model.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.Handle(c, consts.OpAuthRefreshToken, consts.ErrInvalidBody, consts.CodeInvalidBody, richerror.KindInvalid, err)
		return
	}

	if req.RefreshToken == "" {
		richerror.Handle(c, consts.OpAuthRefreshToken, consts.ErrAuthInvalidRefreshToken, consts.CodeInvalidRefreshToken, richerror.KindInvalid, nil)
		return
	}

	meta.Logger.Info("🔄 Refresh token request started")

	claims, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		richerror.HandleWrap(c, consts.OpAuthRefreshToken, consts.ErrAuthInvalidRefreshToken, consts.CodeInvalidRefreshToken, richerror.KindUnauthorized, err)
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		richerror.HandleWrap(c, consts.OpAuthRefreshToken, consts.ErrAuthInvalidUserID, consts.CodeInvalidUserID, richerror.KindInvalid, err)
		return
	}

	accessToken, newRefreshToken, errSent := h.generateTokens(c, consts.OpAuthRefreshToken, userID)
	if errSent {
		meta.Logger.Error("⛔ Token generation failed (refresh)", zap.String("user_id", userID.String()))
		richerror.Handle(c, consts.OpAuthRefreshToken, consts.ErrAuthTokenGenFail, consts.CodeAuthTokenGenFail, richerror.KindInternal, nil)
		return
	}

	meta.Logger.Info("✅ Token refresh successful", zap.String("user_id", userID.String()), zap.Duration("duration", meta.Elapsed()))

	model.SuccessResponse(c, http.StatusOK, model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, consts.MsgTokenRefreshed)
}
