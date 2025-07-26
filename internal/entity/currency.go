package entity

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type Currency struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Code      string         `gorm:"size:10;uniqueIndex;not null" json:"code"` // BTC, ETH, USDT ...
	Name      string         `gorm:"size:100;not null" json:"name"`            // Bitcoin, Ethereum ...
	Precision uint           `gorm:"default:8" json:"precision"`               // تعداد ارقام اعشار مجاز
	IsActive  bool           `gorm:"default:true" json:"is_active"`            // آیا این ارز فعال است
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
