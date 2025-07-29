package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type OrderType string
type OrderSide string
type OrderStatus string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"

	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"

	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCanceled  OrderStatus = "canceled"
	OrderStatusPartial   OrderStatus = "partial" // انجام بخشی از سفارش (قسمت اجرا شده)
	OrderStatusExpired   OrderStatus = "expired"
)

type Order struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID      `gorm:"index;not null"`
	WalletID      uuid.UUID      `gorm:"index;not null"`
	PairID        uuid.UUID      `gorm:"type:uuid;not null"`            // کلید به Pair
	OrderType     OrderType      `gorm:"type:varchar(10);not null"`     // limit, market
	Side          OrderSide      `gorm:"type:varchar(10);not null"`     // buy, sell
	CurrencyPair  string         `gorm:"size:20;not null"`              // BTC/USDT
	Amount        float64        `gorm:"type:decimal(38,18);not null"`  // کل مقدار سفارش
	FilledAmount  float64        `gorm:"type:decimal(38,18);default:0"` // مقدار اجرا شده
	Price         float64        `gorm:"type:decimal(38,18);not null"`  // قیمت سفارش (برای market می‌تونه ۰ باشد)
	Status        OrderStatus    `gorm:"type:varchar(20);default:'pending'"`
	ClientOrderID *string        `gorm:"size:64;index" json:"client_order_id,omitempty"` // شناسه سفارش سمت کلاینت (اختیاری)
	Meta          *string        `gorm:"type:text" json:"meta,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
