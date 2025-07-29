package entity

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type WalletStatus string

const (
	WalletStatusActive   WalletStatus = "active"
	WalletStatusInactive WalletStatus = "inactive"
	WalletStatusFrozen   WalletStatus = "frozen"
	// برای future use: مثلاً suspended, locked, ...
)

type Wallet struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"index;not null" json:"user_id"`
	CurrencyID uuid.UUID `gorm:"index;not null" json:"currency_id"`

	Balance float64 `gorm:"type:decimal(38,18);default:0" json:"balance"` // موجودی قابل برداشت
	Frozen  float64 `gorm:"type:decimal(38,18);default:0" json:"frozen"`  // موجودی بلوکه (سفارش فعال/در انتظار برداشت)
	Total   float64 `gorm:"type:decimal(38,18);default:0" json:"total"`   // مجموع کل موجودی (Balance + Frozen)

	Status WalletStatus `gorm:"type:varchar(16);default:'active';index" json:"status"` // وضعیت فعلی والت

	Meta *string `gorm:"type:text" json:"meta,omitempty"` // اطلاعات اضافی (json قابل توسعه)

	LastActivity *time.Time `json:"last_activity,omitempty"` // آخرین فعالیت

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// --- Relation Optional
	Currency *Currency `gorm:"foreignKey:CurrencyID" json:"currency,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Hook برای تنظیم UUID و تاریخ‌ها هنگام ایجاد
func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New()
	now := time.Now()
	w.CreatedAt = now
	w.UpdatedAt = now
	return
}

func (w *Wallet) BeforeUpdate(tx *gorm.DB) (err error) {
	w.UpdatedAt = time.Now()
	return
}
