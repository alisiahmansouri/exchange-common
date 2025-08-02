package model

import (
	"github.com/google/uuid"
	"time"
)

type WalletActionType string

const (
	ActionDeductFrozen WalletActionType = "DeductFrozen"
	ActionDeposit      WalletActionType = "Deposit"
	// می‌تونی withdrawal, transfer و ... اضافه کنی
)

type WalletAction struct {
	// Unique event ID for traceability and idempotency (مهم برای تکراری نشدن عملیات و مانیتورینگ)
	ActionID  uuid.UUID        `json:"action_id"`
	UserID    uuid.UUID        `json:"user_id"`
	WalletID  uuid.UUID        `json:"wallet_id"`
	Amount    float64          `json:"amount"`
	Action    WalletActionType `json:"action"` // فقط مقادیر مجاز (enum)
	Reason    string           `json:"reason"`
	OrderID   uuid.UUID        `json:"order_id,omitempty"` // Reference to related order, optional but recommended
	PairID    uuid.UUID        `json:"pair_id,omitempty"`  // Reference for trade settlement context
	Ref       string           `json:"ref,omitempty"`      // Any extra reference code (برای ارتباط یا audit یا bridge external systems)
	CreatedAt time.Time        `json:"created_at"`         // زمان ساخت اکشن (برای audit دقیق و ordering)

	// Optional: For distributed tracing/correlation
	TraceID string `json:"trace_id,omitempty"` // ارتباط با سایر سرویس‌ها (مثلاً tracing Jaeger)
	Source  string `json:"source,omitempty"`   // نام سرویس یا subsystem صادرکننده اکشن (مثلاً "settlement-service")
}
