package model

import (
	"github.com/google/uuid"
)

// BulkWalletOp ساختار عملیات گروهی روی کیف پول‌ها (برای ایردراپ، تسویه و ...)
type BulkWalletOp struct {
	UserID     uuid.UUID // کاربر هدف عملیات
	WalletID   uuid.UUID // والت هدف عملیات
	Amount     float64   // مقدار
	OpType     string    // نوع عملیات: deposit/withdraw/freeze/unfreeze/adjust
	Meta       *string   // توضیحات/اطلاعات اضافه (اختیاری)
	OperatorID uuid.UUID // اجراکننده عملیات (مثلاً ادمین)
	Note       *string   // یادداشت اختیاری برای ثبت در لاگ
}
