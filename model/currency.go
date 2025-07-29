package model

import (
	"github.com/alisiahmansouri/exchange-common/entity"
	"github.com/google/uuid"
)

// CurrencyResponse ساختار پاسخ برای اطلاعات یک ارز
type CurrencyResponse struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"` // کد ارز (نماد)
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"` // همان code یا یک نماد نمایشگر (مثل ₿ برای BTC)
	Precision uint      `json:"precision"`
	IsActive  bool      `json:"is_active"`
}

// CurrencyResponseFromEntity تابع تبدیل موجودیت به مدل پاسخ
func CurrencyResponseFromEntity(c entity.Currency) CurrencyResponse {
	return CurrencyResponse{
		ID:        c.ID,
		Code:      c.Code,
		Name:      c.Name,
		Symbol:    c.Code, // یا مقدار سفارشی برای نمایش نماد
		Precision: c.Precision,
		IsActive:  c.IsActive,
	}
}
