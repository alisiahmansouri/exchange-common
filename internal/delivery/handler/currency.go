package handler

import (
	"exchange-common/internal/consts"
	"exchange-common/internal/delivery/requestmeta"
	"exchange-common/internal/model"
	"exchange-common/internal/richerror"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// ListCurrencies godoc
// @Summary      لیست ارزها
// @Description  دریافت لیست تمام ارزهای فعال موجود در سیستم
// @Tags         currency
// @Produce      json
// @Success      200 {object} model.Response[[]model.CurrencyResponse] "لیست ارزها"
// @Failure      500 {object} model.ErrorResponseStruct "خطای داخلی سرور"
// @Router       /currencies [get]
// @Security     BearerAuth
func (h *Handler) ListCurrencies(c *gin.Context) {
	meta := requestmeta.NewRequestMeta(c)
	ctx := c.Request.Context()

	list, err := h.svc.ListActiveCurrencies(ctx)
	if err != nil {
		richerror.HandleWrap(
			c,
			consts.OpCurrencyList,
			consts.ErrCurrencyListFailed,
			consts.CodeCurrencyListError,
			richerror.KindInternal,
			err,
		)
		return
	}

	resp := make([]model.CurrencyResponse, len(list))
	for i, entity := range list {
		resp[i] = model.CurrencyResponseFromEntity(entity)
	}

	meta.Logger.Info("لیست ارزها با موفقیت دریافت شد",
		zap.String("operation", consts.OpCurrencyList),
		zap.Int("count", len(resp)),
		zap.Duration("duration", meta.Elapsed()),
	)

	model.SuccessResponse(c, http.StatusOK, resp)
}
