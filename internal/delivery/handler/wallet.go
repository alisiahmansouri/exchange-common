package handler

import (
	"exchange-common/internal/consts"
	
	"exchange-common/internal/delivery/requestmeta"
	"exchange-common/internal/model"
	"exchange-common/internal/richerror"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ListWallets godoc
// @Summary      لیست کیف پول‌ها
// @Description  دریافت لیست کیف پول‌های یک کاربر
// @Tags         wallet
// @Produce      json
// @Param        user_id  query     string  true  "شناسه کاربر"
// @Success      200 {object} model.Response[[]model.WalletResponse] "لیست کیف پول‌ها"
// @Failure      400 {object} model.ErrorResponseStruct "درخواست نامعتبر"
// @Failure      500 {object} model.ErrorResponseStruct "خطای داخلی سرور"
// @Router       /wallets [get]
// @Security     BearerAuth
func (h *Handler) ListWallets(c *gin.Context) {
	meta := requestmeta.NewRequestMeta(c)
	ctx := c.Request.Context()

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		richerror.HandleWrap(c,.OpWalletCreate, .
		ErrWalletUserInvalidID, .
		CodeInvalidUserID, richerror.KindInvalid, nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		richerror.HandleWrap(c,.OpWalletCreate, .
		ErrWalletUserInvalidID, .
		CodeInvalidUserID, richerror.KindInvalid, err)
		return
	}

	list, err := h.svc.ListWalletsByUser(ctx, userID)
	if err != nil {
		richerror.HandleWrap(c,.OpWalletCreate, .
		ErrWalletListFailed, .
		CodeWalletListError, richerror.KindInternal, err)
		return
	}

	resp := make([]model.WalletResponse, len(list))
	for i, entity := range list {
		resp[i] = model.WalletResponseFromEntity(entity)
	}

	meta.Logger.Info("لیست کیف پول‌ها با موفقیت دریافت شد",
		zap.Int("count", len(resp)),
		zap.Duration("duration", meta.Elapsed()),
	)

	model.SuccessResponse(c, http.StatusOK, resp)
}

// DepositWallet godoc
// @Summary      واریز به کیف پول
// @Description  واریز مبلغ مشخص به کیف پول یک کاربر
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        wallet_id  path      string  true  "شناسه کیف پول"
// @Param        body  body      model.WalletDepositRequest  true  "اطلاعات واریز"
// @Success      200 {object} model.Response[string] "واریز موفق"
// @Failure      400   {object}  model.ErrorResponseStruct "درخواست نامعتبر"
// @Failure      404   {object}  model.ErrorResponseStruct "کیف پول پیدا نشد"
// @Failure      403   {object}  model.ErrorResponseStruct "دسترسی غیرمجاز"
// @Failure      500   {object}  model.ErrorResponseStruct "خطای داخلی سرور"
// @Router       /wallets/deposit [post]
// @Security     BearerAuth
func (h *Handler) DepositWallet(c *gin.Context) {
	meta := requestmeta.NewRequestMeta(c)

	var req model.WalletDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.HandleWrap(c,.OpWalletDeposit, .
		ErrWalletInvalidBody, consts.CodeInvalidRequestBody, richerror.KindInvalid, err)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		richerror.HandleWrap(c,.OpWalletDeposit, .
		ErrWalletUserInvalidID, .
		CodeInvalidUserID, richerror.KindInvalid, err)
		return
	}

	currencyID, err := uuid.Parse(req.CurrencyID)
	if err != nil {
		richerror.HandleWrap(c,.OpWalletDeposit, .
		ErrWalletInvalidID, .
		CodeInvalidCurrencyID, richerror.KindInvalid, err)
		return
	}

	err = h.svc.Deposit(c.Request.Context(), userID, currencyID, req.Amount)
	if err != nil {
		richerror.HandleWrap(c,.OpWalletDeposit, .
		ErrWalletDepositFailed, .
		CodeDepositError, richerror.KindInternal, err)
		return
	}

	meta.Logger.Info("واریز به کیف پول موفق بود",
		zap.String("wallet_id", currencyID.String()),
		zap.String("user_id", userID.String()),
		zap.Float64("amount", req.Amount),
	)

	model.SuccessResponse(c, http.StatusOK,.MsgDepositSuccess)
}

// WithdrawWallet godoc
// @Summary      برداشت از کیف پول
// @Description  برداشت مبلغ مشخص از کیف پول یک کاربر
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        wallet_id  path      string  true  "شناسه کیف پول"
// @Param        body  body      model.WalletWithdrawRequest  true  "اطلاعات برداشت"
// @Success      200 {object} model.Response[string] "برداشت موفق"
// @Failure      400   {object}  model.ErrorResponseStruct "درخواست نامعتبر"
// @Failure      404   {object}  model.ErrorResponseStruct "کیف پول پیدا نشد"
// @Failure      403   {object}  model.ErrorResponseStruct "دسترسی غیرمجاز"
// @Failure      500   {object}  model.ErrorResponseStruct "خطای داخلی سرور"
// @Router       /wallets/{wallet_id}/withdraw [post]
// @Security     BearerAuth
func (h *Handler) WithdrawWallet(c *gin.Context) {
	meta := requestmeta.NewRequestMeta(c)

	walletIDStr := c.Param("wallet_id")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		richerror.HandleWrap(c,.OpWalletWithdraw, .
		ErrWalletInvalidID, .
		CodeInvalidWalletID, richerror.KindInvalid, err)
		return
	}

	var req model.WalletWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		richerror.HandleWrap(c,.OpWalletWithdraw, .
		ErrWalletInvalidBody, consts.CodeInvalidRequestBody, richerror.KindInvalid, err)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		richerror.HandleWrap(c,.OpWalletWithdraw, .
		ErrWalletUserInvalidID, .
		CodeInvalidUserID, richerror.KindInvalid, err)
		return
	}

	err = h.svc.Withdraw(c.Request.Context(), userID, walletID, req.Amount)
	if err != nil {
		richerror.HandleWrap(c,.OpWalletWithdraw, .
		ErrWalletWithdrawFailed, .
		CodeWithdrawError, richerror.KindInternal, err)
		return
	}

	meta.Logger.Info("برداشت از کیف پول موفق بود",
		zap.String("wallet_id", walletID.String()),
		zap.Float64("amount", req.Amount),
	)

	model.SuccessResponse(c, http.StatusOK,.MsgWithdrawSuccess)
}
