package model

import (
	"github.com/google/uuid"
)

// BulkWalletOpType تعریف انواع عملیات قابل پشتیبانی در عملیات گروهی کیف پول
type BulkWalletOpType string

const (
	BulkOpDeposit     BulkWalletOpType = "deposit"
	BulkOpWithdraw    BulkWalletOpType = "withdraw"
	BulkOpFreeze      BulkWalletOpType = "freeze"
	BulkOpUnfreeze    BulkWalletOpType = "unfreeze"
	BulkOpAdjust      BulkWalletOpType = "adjust"
	BulkOpTransferIn  BulkWalletOpType = "transfer-in"
	BulkOpTransferOut BulkWalletOpType = "transfer-out"
)

// BulkWalletOp مدل حرفه‌ای برای یک عملیات گروهی روی کیف پول (برای ایردراپ، اصلاح موجودی، تسویه، ...)
type BulkWalletOp struct {
	UserID     uuid.UUID        `json:"user_id"`        // کاربر هدف عملیات
	WalletID   uuid.UUID        `json:"wallet_id"`      // والت هدف عملیات
	Amount     float64          `json:"amount"`         // مبلغ عملیات (همیشه مثبت)
	OpType     BulkWalletOpType `json:"op_type"`        // نوع عملیات (deposit/withdraw/freeze/unfreeze/adjust)
	Meta       *string          `json:"meta,omitempty"` // اطلاعات اضافه (مثلاً توضیح ایردراپ، شناسه سفارش و ...)
	OperatorID uuid.UUID        `json:"operator_id"`    // اجراکننده عملیات (ادمین/سیستم/ربات)
	Note       *string          `json:"note,omitempty"` // توضیح/یادداشت برای ثبت در لاگ یا تحلیل آینده
}

// ExampleBulkOps یک نمونه لیست عملیات گروهی برای ایردراپ یا تسویه
var ExampleBulkOps = []BulkWalletOp{
	{
		UserID:     uuid.New(),
		WalletID:   uuid.New(),
		Amount:     100.0,
		OpType:     BulkOpDeposit,
		Meta:       ptr("airdrop event #313"),
		OperatorID: uuid.New(),
		Note:       ptr("ایردراپ اسفند ۱۴۰۳"),
	},
	{
		UserID:     uuid.New(),
		WalletID:   uuid.New(),
		Amount:     250.0,
		OpType:     BulkOpFreeze,
		Meta:       ptr("freeze for big order 821"),
		OperatorID: uuid.New(),
		Note:       ptr("سفارش VIP فریز"),
	},
}

// ptr یک کمک ساده برای مقداردهی اشاره‌گر به string
func ptr(s string) *string { return &s }

func (op BulkWalletOp) NoteString() string {
	if op.Note != nil {
		return *op.Note
	}
	return ""
}
