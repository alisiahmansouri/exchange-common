package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Pair struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BaseCurrencyID  uuid.UUID `gorm:"type:uuid;not null"`
	QuoteCurrencyID uuid.UUID `gorm:"type:uuid;not null"`

	BaseCurrency  Currency `gorm:"foreignKey:BaseCurrencyID;references:ID"`
	QuoteCurrency Currency `gorm:"foreignKey:QuoteCurrencyID;references:ID"`

	PricePrecision  uint
	AmountPrecision uint

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
