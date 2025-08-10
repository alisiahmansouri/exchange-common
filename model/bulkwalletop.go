package model

import (
	"strings"

	"github.com/google/uuid"
)

// BulkWalletOpType تعریف انواع عملیات قابل پشتیبانی در عملیات گروهی کیف پول
type BulkWalletOpType string

const (
	BulkOpDeposit      BulkWalletOpType = "deposit"
	BulkOpWithdraw     BulkWalletOpType = "withdraw"
	BulkOpFreeze       BulkWalletOpType = "freeze"
	BulkOpUnfreeze     BulkWalletOpType = "unfreeze"
	BulkOpAdjust       BulkWalletOpType = "adjust"
	BulkOpTransferIn   BulkWalletOpType = "transfer-in"  // فقط برای گزارش/لاگ
	BulkOpTransferOut  BulkWalletOpType = "transfer-out" // فقط برای گزارش/لاگ
	BulkOpDeductFrozen BulkWalletOpType = "deduct-frozen"
)

// BulkWalletOp مدل حرفه‌ای برای یک عملیات گروهی روی کیف پول (ایردراپ، اصلاح موجودی، تسویه، ...)
// نکته درباره Amount:
//   - برای deposit/withdraw/freeze/unfreeze/deduct-frozen باید > 0 باشد.
//   - برای adjust می‌تواند مثبت یا منفی باشد (صفر مجاز نیست).
//   - در UseCase، اعتبارسنجی دقیق‌تر اعمال می‌شود.
type BulkWalletOp struct {
	OpType     BulkWalletOpType `json:"op_type"`               // نوع عملیات
	UserID     uuid.UUID        `json:"user_id"`               // کاربر هدف
	WalletID   uuid.UUID        `json:"wallet_id,omitempty"`   // کیف‌پول هدف؛ برای deposit می‌تواند خالی باشد
	CurrencyID uuid.UUID        `json:"currency_id,omitempty"` // فقط وقتی WalletID خالی است و op=deposit
	Amount     float64          `json:"amount"`                // مقدار
	Meta       *string          `json:"meta,omitempty"`        // متادیتا (JSON string یا متن)
	OperatorID uuid.UUID        `json:"operator_id,omitempty"` // اجراکننده (ادمین/سیستم) — کنترل سطح دسترسی بیرون
	Note       *string          `json:"note,omitempty"`        // توضیح کوتاه برای لاگ/تحلیل
}

// NoteString خروجی تمیز از Note می‌دهد
func (op BulkWalletOp) NoteString() string {
	if op.Note == nil {
		return ""
	}
	return strings.TrimSpace(*op.Note)
}

// ExampleBulkOps یک نمونه لیست عملیات گروهی برای ایردراپ/تسویه
var ExampleBulkOps = []BulkWalletOp{
	{
		OpType: BulkOpDeposit,
		UserID: uuid.New(),
		// WalletID خالی است؛ با CurrencyID ساخته/یافت می‌شود
		CurrencyID: uuid.New(),
		Amount:     100.0,
		Meta:       ptr("airdrop event #313"),
		OperatorID: uuid.New(),
		Note:       ptr("ایردراپ اسفند ۱۴۰۳"),
	},
	{
		OpType:     BulkOpFreeze,
		UserID:     uuid.New(),
		WalletID:   uuid.New(),
		Amount:     250.0,
		Meta:       ptr("freeze for big order 821"),
		OperatorID: uuid.New(),
		Note:       ptr("سفارش VIP فریز"),
	},
	{
		OpType:     BulkOpAdjust,
		UserID:     uuid.New(),
		WalletID:   uuid.New(),
		Amount:     -50.0, // کاهش دستی (نمونه)
		Meta:       ptr("manual reconciliation"),
		OperatorID: uuid.New(),
		Note:       ptr("تنظیم اختلاف موجودی"),
	},
}

// ptr یک کمک ساده برای مقداردهی اشاره‌گر به string
func ptr(s string) *string { return &s }
