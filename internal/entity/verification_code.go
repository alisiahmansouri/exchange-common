package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type VerificationCode struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	Identifier string         `gorm:"index;not null" json:"identifier"`
	HashedCode string         `gorm:"not null" json:"hashed-code"`
	ExpiresAt  time.Time      `gorm:"not null" json:"expires_at"`
	IsUsed     bool           `gorm:"default:false" json:"is_used"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Purpose    string         `gorm:"type:varchar(32);index" json:"purpose"`
	Channel    string         `gorm:"type:varchar(20);index" json:"channel"`
}

func (v *VerificationCode) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	v.CreatedAt = time.Now()
	return
}
