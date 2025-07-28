package entity

import (
	"time"

	"gorm.io/gorm"
)

type OrderType string
type OrderStatus string

const (
	OrderTypeBuy  OrderType = "buy"
	OrderTypeSell OrderType = "sell"

	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCanceled  OrderStatus = "canceled"
)

type Order struct {
	ID           uint        `gorm:"primaryKey;autoIncrement"`
	UserID       uint        `gorm:"index;not null"`
	WalletID     uint        `gorm:"index;not null"` // کیف پول مرتبط
	OrderType    OrderType   `gorm:"type:varchar(10);not null"`
	CurrencyPair string      `gorm:"size:20;not null"` // مثلا BTC/USDT
	Amount       float64     `gorm:"type:decimal(30,8);not null"`
	Price        float64     `gorm:"type:decimal(30,8);not null"` // قیمت سفارش
	Status       OrderStatus `gorm:"type:varchar(20);default:'pending'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
