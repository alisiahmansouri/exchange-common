package entity

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type Currency struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Code      string         `gorm:"size:10;uniqueIndex;not null" json:"code"`       // BTC, ETH, USDT ...
	Name      string         `gorm:"size:100;not null" json:"name"`                  // Bitcoin, Ethereum ...
	Symbol    string         `gorm:"size:10;" json:"symbol,omitempty"`               // ₿, Ξ, $
	Type      string         `gorm:"size:20;default:'crypto'" json:"type,omitempty"` // crypto, fiat, token
	Chain     string         `gorm:"size:30;" json:"chain,omitempty"`                // bitcoin, ethereum, tron, etc.
	Meta      *string        `gorm:"type:text" json:"meta,omitempty"`                // اطلاعات اضافی (برای توسعه آینده، مثلا شبکه‌ها و...)
	Precision uint           `gorm:"default:8" json:"precision"`                     // تعداد اعشار مجاز
	IsActive  bool           `gorm:"default:true" json:"is_active"`                  // آیا فعال است
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Currency) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	return
}

func (c *Currency) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}
