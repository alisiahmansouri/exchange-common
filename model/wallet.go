package model

import (
	"github.com/alisiahmansouri/exchange-common/entity"
	"github.com/google/uuid"
)

// ---- درخواست ایجاد کیف پول جدید ----
type WalletCreateRequest struct {
	UserID     string  `json:"user_id" binding:"required,uuid4"`
	CurrencyID string  `json:"currency_id" binding:"required,uuid4"`
	Balance    float64 `json:"balance" binding:"gte=0"`
	Status     string  `json:"status" binding:"omitempty,oneof=active inactive frozen"` // به طور پیش‌فرض active
	Meta       *string `json:"meta,omitempty"`
}

// تبدیل به entity.Wallet حرفه‌ای
func (r *WalletCreateRequest) ToEntity() entity.Wallet {
	return entity.Wallet{
		ID:         uuid.New(),
		UserID:     mustParseUUID(r.UserID),
		CurrencyID: mustParseUUID(r.CurrencyID),
		Balance:    r.Balance,
		Total:      r.Balance, // در ایجاد اولیه معمولاً Total=Balance، مگر حالت خاص
		Frozen:     0,
		Status:     entity.WalletStatus(r.Status),
		Meta:       r.Meta,
	}
}

// ---- درخواست واریز ----
type WalletDepositRequest struct {
	UserID     string  `json:"user_id" binding:"required,uuid4"`
	CurrencyID string  `json:"currency_id" binding:"required,uuid4"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	Meta       *string `json:"meta,omitempty"`
}

// ---- درخواست برداشت ----
type WalletWithdrawRequest struct {
	UserID     string  `json:"user_id" binding:"required,uuid4"`
	CurrencyID string  `json:"currency_id" binding:"required,uuid4"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	Meta       *string `json:"meta,omitempty"`
}

// ---- درخواست انتقال داخلی بین کیف پول‌ها ----
type WalletTransferRequest struct {
	UserID         string  `json:"user_id" binding:"required,uuid4"`
	FromCurrencyID string  `json:"from_currency_id" binding:"required,uuid4"`
	ToCurrencyID   string  `json:"to_currency_id" binding:"required,uuid4"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Meta           *string `json:"meta,omitempty"`
}

// ---- مدل پاسخ کیف پول ----
type WalletResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	CurrencyID   string  `json:"currency_id"`
	Balance      float64 `json:"balance"`
	Frozen       float64 `json:"frozen"`
	Total        float64 `json:"total"`
	Status       string  `json:"status"`
	Meta         *string `json:"meta,omitempty"`
	LastActivity *string `json:"last_activity,omitempty"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// ---- تبدیل entity.Wallet به WalletResponse ----
func WalletResponseFromEntity(wallet entity.Wallet) WalletResponse {
	var lastActivity *string
	if wallet.LastActivity != nil {
		str := wallet.LastActivity.Format("2006-01-02T15:04:05Z07:00")
		lastActivity = &str
	}
	return WalletResponse{
		ID:           wallet.ID.String(),
		UserID:       wallet.UserID.String(),
		CurrencyID:   wallet.CurrencyID.String(),
		Balance:      wallet.Balance,
		Frozen:       wallet.Frozen,
		Total:        wallet.Total,
		Status:       string(wallet.Status),
		Meta:         wallet.Meta,
		LastActivity: lastActivity,
		CreatedAt:    wallet.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    wallet.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ---- ابزار کمکی ----
func mustParseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
