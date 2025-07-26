package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email           *string        `gorm:"uniqueIndex;null" json:"email,omitempty"` // nullable
	Phone           *string        `gorm:"uniqueIndex;null" json:"phone,omitempty"` // nullable
	PasswordHash    string         `gorm:"not null" json:"-"`
	FullName        string         `json:"full_name"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	IsEmailVerified bool           `gorm:"default:false" json:"is_email_verified"`
	IsPhoneVerified bool           `gorm:"default:false" json:"is_phone_verified"`
	TwoFAEnabled    bool           `gorm:"default:false" json:"two_fa_enabled"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}
