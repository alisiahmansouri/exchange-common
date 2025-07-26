package entity

import (
	"time"

	"gorm.io/gorm"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID           uint              `gorm:"primaryKey;autoIncrement"`
	OrderID      uint              `gorm:"index;not null"`
	FromWalletID uint              `gorm:"index;not null"` // کیف پول فرستنده
	ToWalletID   uint              `gorm:"index;not null"` // کیف پول گیرنده
	Amount       float64           `gorm:"type:decimal(30,8);not null"`
	Fee          float64           `gorm:"type:decimal(30,8);default:0"` // کارمزد تراکنش
	Status       TransactionStatus `gorm:"type:varchar(20);default:'pending'"`
	TxHash       string            `gorm:"size:100;uniqueIndex"` // هش تراکنش بلاکچین (اختیاری)
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
