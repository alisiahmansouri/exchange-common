package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// انواع تراکنش
type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"  // واریز
	TransactionTypeWithdraw TransactionType = "withdraw" // برداشت
	TransactionTypeTransfer TransactionType = "transfer" // انتقال داخلی بین والت‌ها
	TransactionTypeTrade    TransactionType = "trade"    // تسویه سفارش بازار
	TransactionTypeFee      TransactionType = "fee"      // کارمزد (درصورت نیاز)
)

// وضعیت‌های تراکنش
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCanceled  TransactionStatus = "canceled"
)

type Transaction struct {
	ID           uuid.UUID         `gorm:"type:uuid;primaryKey"`                            // شناسه یکتا
	OrderID      *uuid.UUID        `gorm:"type:uuid;index" json:"order_id,omitempty"`       // Nullable: تراکنش‌های غیرمارکتی Order ندارند
	FromWalletID *uuid.UUID        `gorm:"type:uuid;index" json:"from_wallet_id,omitempty"` // Nullable: واریز از بیرون (والت ندارد)
	ToWalletID   *uuid.UUID        `gorm:"type:uuid;index" json:"to_wallet_id,omitempty"`   // Nullable: برداشت به بیرون (والت ندارد)
	Type         TransactionType   `gorm:"type:varchar(16);not null;default:'transfer'" json:"type"`
	Amount       float64           `gorm:"type:decimal(38,18);not null"`                  // مقدار تراکنش
	Fee          float64           `gorm:"type:decimal(38,18);default:0"`                 // کارمزد
	Status       TransactionStatus `gorm:"type:varchar(20);default:'pending'"`            // وضعیت
	TxHash       *string           `gorm:"size:100;uniqueIndex" json:"tx_hash,omitempty"` // هش بلاکچین (برای واریز/برداشت)
	Note         *string           `gorm:"type:text" json:"note,omitempty"`               // توضیحات اضافی (مثلاً علت خطا)
	Meta         *string           `gorm:"type:text" json:"meta,omitempty"`               // future-proof (برای هرنوع اطلاعات دلخواه)
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `gorm:"index" json:"-"`

	// ارتباطات (اختیاری برای GORM preload)
	Order      *Order  `gorm:"foreignKey:OrderID"`
	FromWallet *Wallet `gorm:"foreignKey:FromWalletID"`
	ToWallet   *Wallet `gorm:"foreignKey:ToWalletID"`
}
