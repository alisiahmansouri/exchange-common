package entity

import (
	"github.com/google/uuid"
	"time"
)

// نوع عملیات (Deposit, Withdraw, Freeze, Unfreeze, Transfer, Adjust, ...)
type WalletLogType string

const (
	WalletLogDeposit      WalletLogType = "deposit"
	WalletLogWithdraw     WalletLogType = "withdraw"
	WalletLogTransferIn   WalletLogType = "transfer-in"
	WalletLogTransferOut  WalletLogType = "transfer-out"
	WalletLogFreeze       WalletLogType = "freeze"
	WalletLogUnfreeze     WalletLogType = "unfreeze"
	WalletLogDeductFrozen WalletLogType = "deduct-frozen"
	WalletLogChangeStatus WalletLogType = "change-status"
	WalletLogAdjust       WalletLogType = "admin-adjust"
)

type WalletLog struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey"`
	WalletID  uuid.UUID     `gorm:"type:uuid;index;not null"`
	UserID    uuid.UUID     `gorm:"type:uuid;index;not null"`
	LogType   WalletLogType `gorm:"type:varchar(32);index;not null"`
	Amount    float64       `gorm:"type:decimal(38,18);not null"`
	Meta      *string       `gorm:"type:text" json:"meta,omitempty"`
	CreatedAt time.Time     `gorm:"not null"`
}

func NewWalletLog(walletID, userID uuid.UUID, logType string, amount float64, meta *string) WalletLog {
	return WalletLog{
		ID:        uuid.New(),
		WalletID:  walletID,
		UserID:    userID,
		LogType:   WalletLogType(logType),
		Amount:    amount,
		Meta:      meta,
		CreatedAt: time.Now(),
	}
}
