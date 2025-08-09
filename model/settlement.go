package model

import (
	"github.com/google/uuid"
	"time"
)

const (
	// نسخه‌ی فعلی پیام رویداد؛ برای سازگاری رو به جلو/عقب
	SettleTradeEventVersion = 1
)

// SettleTradeEvent is the event sent from matching engine to settlement service for trade settlement.
type SettleTradeEvent struct {
	EventID       uuid.UUID `json:"event_id"`                 // For idempotency & audit (UUIDv5 recommended)
	Version       int       `json:"version"`                  // For backward/forward compatibility (use SettleTradeEventVersion)
	PairID        uuid.UUID `json:"pair_id"`                  // Trading pair ID (for sharding & ordering)
	Sequence      uint64    `json:"sequence"`                 // Monotonic per pair (ordering)
	TakerOrderID  uuid.UUID `json:"taker_order_id"`           // Incoming (taker) order ID
	MakerOrderID  uuid.UUID `json:"maker_order_id"`           // Book (maker) order ID
	MatchAmount   float64   `json:"match_amount"`             // Matched base amount
	TradePrice    float64   `json:"trade_price"`              // Matched price
	TraceID       string    `json:"trace_id,omitempty"`       // For distributed tracing (optional)
	CorrelationID string    `json:"correlation_id,omitempty"` // For cross-service tracing (optional)
	CreatedAt     time.Time `json:"created_at"`               // UTC time of event creation
	// Future fields:
	// Metadata map[string]string `json:"metadata,omitempty"`
}

// نتیجه ذخیره‌شده تسویه (در Redis)
type SettleTradeResult struct {
	Status       string    `json:"status"`
	SettlementID string    `json:"settlement_id,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	Message      string    `json:"message,omitempty"`
}
