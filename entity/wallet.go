package entity

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"index;not null" json:"user_id"`
	CurrencyID uuid.UUID      `gorm:"index;not null" json:"currency_id"`
	Balance    float64        `gorm:"type:decimal(30,8);default:0" json:"balance"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// Hook برای تنظیم UUID و تاریخ‌ها هنگام ایجاد
func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New()
	now := time.Now()
	w.CreatedAt = now
	w.UpdatedAt = now
	return
}

// Hook برای بروزرسانی زمان
func (w *Wallet) BeforeUpdate(tx *gorm.DB) (err error) {
	w.UpdatedAt = time.Now()
	return
}
