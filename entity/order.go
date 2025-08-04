package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// --- Enumerations ---
type OrderType string
type OrderSide string
type OrderStatus string
type OrderTimeInForce string // برای سفارشات پیشرفته (GTC, IOC, FOK, ...)

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"

	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"

	OrderStatusPending   OrderStatus = "pending"   // در انتظار ثبت/فعال‌سازی
	OrderStatusActive    OrderStatus = "active"    // فعال در دفتر سفارشات
	OrderStatusPartial   OrderStatus = "partial"   // بخشی از سفارش اجرا شده (باقیمانده فعال)
	OrderStatusCompleted OrderStatus = "completed" // تمام سفارش اجرا شد
	OrderStatusCanceled  OrderStatus = "canceled"  // لغو شده توسط کاربر/سیستم
	OrderStatusExpired   OrderStatus = "expired"   // منقضی (مثلاً به علت TimeInForce)
	OrderStatusRejected  OrderStatus = "rejected"  // رد شده (مثلاً خطای اعتبارسنجی/محدودیت)
)

const (
	// Time in Force انواع سفارشات حرفه‌ای:
	TimeInForceGTC OrderTimeInForce = "GTC" // Good Till Cancel
	TimeInForceIOC OrderTimeInForce = "IOC" // Immediate or Cancel
	TimeInForceFOK OrderTimeInForce = "FOK" // Fill or Kill
)

// --- Entity: Order ---
type Order struct {
	ID            uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	WalletID      uuid.UUID        `gorm:"type:uuid;not null;index" json:"wallet_id"`
	PairID        uuid.UUID        `gorm:"type:uuid;not null;index" json:"pair_id"`
	SettlementID  *uuid.UUID       `gorm:"type:uuid;index" json:"settlement_id,omitempty"`              //
	OrderType     OrderType        `gorm:"type:varchar(10);not null" json:"order_type"`                 // limit, market
	Side          OrderSide        `gorm:"type:varchar(10);not null" json:"side"`                       // buy, sell
	Amount        float64          `gorm:"type:decimal(38,18);not null" json:"amount"`                  // کل مقدار سفارش
	FilledAmount  float64          `gorm:"type:decimal(38,18);not null;default:0" json:"filled_amount"` // مقدار اجرا شده
	Price         float64          `gorm:"type:decimal(38,18);not null" json:"price"`                   // قیمت سفارش (برای market اختیاری/۰)
	Status        OrderStatus      `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	TimeInForce   OrderTimeInForce `gorm:"type:varchar(8);default:'GTC'" json:"time_in_force"` // future: GTC/IOC/FOK
	ClientOrderID *string          `gorm:"size:64;index" json:"client_order_id,omitempty"`     // شناسه سمت کلاینت (برای تطبیق سریع)
	Meta          *string          `gorm:"type:text" json:"meta,omitempty"`                    // json, برای ثبت مقادیر اضافی (fee, device, ip, ...)
	ExecutedAt    *time.Time       `json:"executed_at,omitempty"`                              // زمان اجرای کامل (settlement)
	ExpiresAt     *time.Time       `json:"expires_at,omitempty"`                               // اگر سفارش زمان انقضا دارد (for IOC/FOK)

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// --- GORM Hooks برای audit/ایمنی بیشتر ---
func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now
	o.Status = OrderStatusPending // تضمین اولیه بودن وضعیت
	if o.FilledAmount < 0 {
		o.FilledAmount = 0
	}
	return nil
}

func (o *Order) BeforeUpdate(tx *gorm.DB) (err error) {
	o.UpdatedAt = time.Now()
	return nil
}

// --- Index پیشنهادی برای performance ---
// CREATE INDEX idx_orders_pair_status ON orders(pair_id, status);
// CREATE INDEX idx_orders_user_status ON orders(user_id, status);

// --- Future extensions ---
// - FeePercent, FeeFixed
// - Refs to related Transaction(s) or Trade(s) if multi-match
// - OCO (one cancels the other), PostOnly, ReduceOnly
// - OrderSource (web, api, mobile), etc.
