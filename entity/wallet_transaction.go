package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// انواع تراکنش کیف پول
type WalletTransactionType string

const (
	WalletTxnDeposit          WalletTransactionType = "deposit"           // واریز مستقیم
	WalletTxnWithdraw         WalletTransactionType = "withdraw"          // برداشت مستقیم
	WalletTxnFreeze           WalletTransactionType = "freeze"            // فریز کردن مبلغ (برای سفارش)
	WalletTxnUnfreeze         WalletTransactionType = "unfreeze"          // آزادسازی مبلغ فریز شده
	WalletTxnDeductFrozen     WalletTransactionType = "deduct_frozen"     // کسر مبلغ فریز شده (تسویه سفارش)
	WalletTxnInternalTransfer WalletTransactionType = "internal_transfer" // انتقال داخلی (بین دو والت)
	WalletTxnTrade            WalletTransactionType = "trade"             // تسویه خرید/فروش بازار
	WalletTxnFee              WalletTransactionType = "fee"               // کارمزد (کسر از والت)
)

// وضعیت تراکنش کیف پول (برای audit دقیق‌تر)
type WalletTransactionStatus string

const (
	WalletTxnStatusPending   WalletTransactionStatus = "pending"
	WalletTxnStatusCompleted WalletTransactionStatus = "completed"
	WalletTxnStatusFailed    WalletTransactionStatus = "failed"
	WalletTxnStatusCanceled  WalletTransactionStatus = "canceled"
)

type WalletTransaction struct {
	ID            uuid.UUID               `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID               `gorm:"type:uuid;not null;index"`
	WalletID      uuid.UUID               `gorm:"type:uuid;not null;index"`
	Type          WalletTransactionType   `gorm:"type:varchar(24);not null;index"`
	Status        WalletTransactionStatus `gorm:"type:varchar(16);not null;default:'completed';index"`
	Amount        float64                 `gorm:"type:decimal(38,18);not null"` // مثبت یا منفی بسته به نوع
	Fee           float64                 `gorm:"type:decimal(38,18);default:0"`
	BalanceBefore float64                 `gorm:"type:decimal(38,18);not null"` // برای audit: مانده قبل
	BalanceAfter  float64                 `gorm:"type:decimal(38,18);not null"` // برای audit: مانده بعد

	TransactionID *uuid.UUID `gorm:"type:uuid;index"`  // ارجاع به تراکنش اصلی سیستم (nullable)
	OrderID       *uuid.UUID `gorm:"type:uuid;index"`  // ارجاع به سفارش (nullable)
	RefID         *uuid.UUID `gorm:"type:uuid;index"`  // ارجاع به هر entity دیگر (مثلاً تسویه کارمزد)
	RefType       *string    `gorm:"type:varchar(32)"` // نوع entity مرجع (مثلاً order, withdrawal و ...)
	Meta          *string    `gorm:"type:text"`        // اطلاعات جانبی: IP, Device, توضیحات و ...

	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
