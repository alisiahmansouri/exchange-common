package entity

import (
	"fmt"
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
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index:idx_user_wallet_created_at,priority:1"`
	WalletID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_wallet_created_at,priority:2"`

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

	CreatedAt time.Time      `gorm:"not null;index:idx_user_wallet_created_at,priority:3"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *WalletTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	// اگر ID مقداردهی نشده، اتوماتیک بساز
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = now
	}
	// قرارداد: Amount نباید صفر یا منفی باشد (مگر برای Typeهای خاص)
	if t.Amount <= 0 && t.Type != WalletTxnFee { // فقط برای Fee مقدار منفی یا صفر مجاز است
		return fmt.Errorf("amount must be positive except for fee")
	}
	// قرارداد: Fee نباید منفی باشد
	if t.Fee < 0 {
		return fmt.Errorf("fee cannot be negative")
	}
	// قرارداد: BalanceAfter باید برابر با BalanceBefore + Amount - Fee باشد
	calculated := t.BalanceBefore + t.Amount - t.Fee
	if t.BalanceAfter != calculated {
		t.BalanceAfter = calculated // یا خطا برگردان، بسته به سیاست پروژه
		// return fmt.Errorf("balanceAfter is invalid")
	}
	// اگر Status خالی است، مقداردهی کن
	if t.Status == "" {
		t.Status = WalletTxnStatusCompleted
	}
	return nil
}

func (t *WalletTransaction) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now()
	// قوانین دیگر همانند BeforeCreate
	return nil
}
