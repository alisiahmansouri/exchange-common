package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Pair struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	BaseCurrencyID  uuid.UUID `gorm:"type:uuid;not null"`
	QuoteCurrencyID uuid.UUID `gorm:"type:uuid;not null"`

	BaseCurrency  Currency `gorm:"foreignKey:BaseCurrencyID;references:ID"`
	QuoteCurrency Currency `gorm:"foreignKey:QuoteCurrencyID;references:ID"`

	Symbol          string  `gorm:"size:20;uniqueIndex;not null"`                          // مثلا BTCUSDT یا ETHUSDT
	PricePrecision  uint    `gorm:"not null"`                                              // دقت قیمت (مثلا 2 رقم اعشار)
	AmountPrecision uint    `gorm:"not null"`                                              // دقت مقدار (مثلا 6 رقم اعشار)
	MinOrderAmount  float64 `gorm:"type:decimal(38,18);default:0" json:"min_order_amount"` // حداقل مقدار مجاز سفارش
	MaxOrderAmount  float64 `gorm:"type:decimal(38,18);default:0" json:"max_order_amount"` // حداکثر مقدار مجاز سفارش (اختیاری)
	IsActive        bool    `gorm:"default:true" json:"is_active"`
	Meta            *string `gorm:"type:text" json:"meta,omitempty"` // اطلاعات اضافی (مثل فیلد fee، config و ...)

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
