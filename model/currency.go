package model

import (
	"github.com/alisiahmansouri/exchange-common/entity"
	"github.com/google/uuid"
)

// CurrencyResponse ساختار پاسخ برای اطلاعات یک ارز
type CurrencyResponse struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol,omitempty"` // نماد واقعی، مثلا ₿ برای BTC
	Type      string    `json:"type,omitempty"`   // crypto, fiat, token
	Chain     string    `json:"chain,omitempty"`  // مثلا bitcoin, ethereum, tron
	Meta      *string   `json:"meta,omitempty"`   // اطلاعات اضافی
	Precision uint      `json:"precision"`
	IsActive  bool      `json:"is_active"`
}

// CurrencyResponseFromEntity تبدیل entity به مدل پاسخ حرفه‌ای
func CurrencyResponseFromEntity(c entity.Currency) CurrencyResponse {
	return CurrencyResponse{
		ID:        c.ID,
		Code:      c.Code,
		Name:      c.Name,
		Symbol:    c.Symbol, // نماد واقعی (اختیاری!)
		Type:      c.Type,
		Chain:     c.Chain,
		Meta:      c.Meta,
		Precision: c.Precision,
		IsActive:  c.IsActive,
	}
}
