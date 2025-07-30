package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Trade struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PairID uuid.UUID `gorm:"type:uuid;not null;index" json:"pair_id"`    // جفت ارز
	Price  float64   `gorm:"type:decimal(38,18);not null" json:"price"`  // قیمت معامله
	Amount float64   `gorm:"type:decimal(38,18);not null" json:"amount"` // مقدار معامله

	TakerOrderID uuid.UUID `gorm:"type:uuid;not null;index" json:"taker_order_id"` // سفارش فعال (همیشه کاربر جدید)
	MakerOrderID uuid.UUID `gorm:"type:uuid;not null;index" json:"maker_order_id"` // سفارش سمت مقابل (در book بوده)
	TakerUserID  uuid.UUID `gorm:"type:uuid;not null;index" json:"taker_user_id"`  // کاربر taker (همیشه سفارش جدید)
	MakerUserID  uuid.UUID `gorm:"type:uuid;not null;index" json:"maker_user_id"`  // کاربر maker (سمت book)

	// کارمزدها (می‌تواند بسته به بازار برای دو طرف متفاوت باشد)
	TakerFee float64 `gorm:"type:decimal(38,18);default:0" json:"taker_fee"`
	MakerFee float64 `gorm:"type:decimal(38,18);default:0" json:"maker_fee"`

	// optional: ثبت meta یا توضیح
	Meta *string `gorm:"type:text" json:"meta,omitempty"`

	// زمان‌ها
	CreatedAt time.Time  `gorm:"not null;index" json:"created_at"`
	SettledAt *time.Time `json:"settled_at,omitempty"` // اگر کامل settle شده باشد

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Trade) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	return nil
}
