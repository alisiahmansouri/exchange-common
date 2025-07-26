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

// Logout godoc
// @Summary خروج کاربر
// @Description باطل کردن توکن دسترسی و پایان جلسه کاربر
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} model.Response[model.SimpleMessageResponse] "خروج موفق"
// @Failure 400 {object} model.ErrorResponseStruct "درخواست نادرست یا هدر احراز هویت ناموجود"
// @Failure 401 {object} model.ErrorResponseStruct "توکن نامعتبر یا کاربر احراز نشده"
// @Failure 500 {object} model.ErrorResponseStruct "خطای داخلی سرور"
// @Router /auth/logout [post]
// @Security BearerAuth
func (h *Handler) Logout(c *gin.Context) {
	meta := requestmeta.NewRequestMeta(c)

	// استخراج userID از context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		richerror.Handle(c, consts.OpAuthLogout, consts.ErrAuthNoAuthHeader, consts.CodeAuthHeaderMissing, richerror.KindUnauthorized, nil)
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		richerror.Handle(c, consts.OpAuthLogout, consts.ErrAuthUserIDTypeInvalid, consts.CodeUserIDTypeInvalid, richerror.KindInternal, nil)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		richerror.Handle(c, consts.OpAuthLogout, consts.ErrAuthInvalidUserID, consts.CodeInvalidUserID, richerror.KindInvalid, err)
		return
	}

	// استخراج توکن از context
	tokenVal, exists := c.Get("jwt_token")
	if !exists {
		richerror.Handle(c, consts.OpAuthLogout, consts.ErrAuthTokenNotFound, consts.CodeTokenNotFoundInContext, richerror.KindUnauthorized, nil)
		return
	}
	token, ok := tokenVal.(string)
	if !ok || token == "" {
		richerror.Handle(c, consts.OpAuthLogout, consts.ErrAuthEmptyToken, consts.CodeTokenEmpty, richerror.KindInvalid, nil)
		return
	}

	meta.Logger.Info("🔒 Logout request started", zap.String("user_id", userID.String()))

	if err := h.jwtService.RevokeToken(token); err != nil {
		richerror.HandleWrap(c, consts.OpAuthLogout, consts.ErrAuthRevokeFail, consts.CodeJWTRevokeError, richerror.KindInternal, err)
		return
	}

	meta.Logger.Info("✅ Logout successful", zap.String("user_id", userID.String()), zap.Duration("duration", meta.Elapsed()))

	model.SuccessResponse(c, http.StatusOK, model.SimpleMessageResponse{
		Success: true,
		Message: consts.MsgLogoutSuccess,
	})
}
