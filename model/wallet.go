package model

import (
	"github.com/alisiahmansouri/exchange-common/entity"
	"github.com/google/uuid"
)

// WalletCreateRequest ساختار درخواست ایجاد کیف پول جدید
type WalletCreateRequest struct {
	UserID     string  `json:"user_id" binding:"required,uuid4"`
	CurrencyID string  `json:"currency_id" binding:"required,uuid4"`
	Balance    float64 `json:"balance" binding:"gte=0"` // موجودی اولیه باید صفر یا مثبت باشد
}

// تبدیل WalletCreateRequest به entity.Wallet
func (r *WalletCreateRequest) ToEntity() entity.Wallet {
	return entity.Wallet{
		ID:         uuid.New(),
		UserID:     mustParseUUID(r.UserID),
		CurrencyID: mustParseUUID(r.CurrencyID),
		Balance:    r.Balance,
	}
}

// WalletDepositRequest ساختار درخواست واریز به کیف پول
type WalletDepositRequest struct {
	UserID     string  `json:"user_id" binding:"required,uuid4"`
	CurrencyID string  `json:"currency_id" binding:"required,uuid4"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
}

// WalletWithdrawRequest ساختار درخواست برداشت از کیف پول
type WalletWithdrawRequest struct {
	UserID string  `json:"user_id" binding:"required,uuid4"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// کمکی برای parse UUID با فرض اینکه ورودی همیشه درست باشد (در هندلر اعتبارسنجی انجام شده)
func mustParseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}

// مدل پاسخ WalletResponse برای برگشت به کلاینت
type WalletResponse struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	CurrencyID string  `json:"currency_id"`
	Balance    float64 `json:"balance"`
}

// تبدیل entity.Wallet به WalletResponse
func WalletResponseFromEntity(wallet entity.Wallet) WalletResponse {
	return WalletResponse{
		ID:         wallet.ID.String(),
		UserID:     wallet.UserID.String(),
		CurrencyID: wallet.CurrencyID.String(),
		Balance:    wallet.Balance,
	}
}
