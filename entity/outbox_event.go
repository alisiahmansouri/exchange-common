package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

const (
	OutboxStatusPending = "pending"
	OutboxStatusSent    = "sent"
	OutboxStatusError   = "error"
)

type OutboxEvent struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	EventType    string    `gorm:"size:64"`
	Payload      []byte    `gorm:"type:jsonb"`
	Status       string    `gorm:"size:16;default:pending"`
	CreatedAt    time.Time
	SentAt       *time.Time
	ErrorMessage *string `gorm:"size:255"` // اختیاری: ثبت آخرین خطای ارسال
}

// Optional: Custom table name (if you want)
func (OutboxEvent) TableName() string { return "match_events_outbox" }

// Optional: Unmarshal payload helper
func (e *OutboxEvent) UnmarshalPayload(v interface{}) error {
	return json.Unmarshal(e.Payload, v)
}
