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
	EventType    string    `gorm:"size:64;not null;index:idx_outbox_type_dedup,priority:1"`
	DedupKey     string    `gorm:"type:char(64);not null;index:idx_outbox_type_dedup,unique,priority:2"` // SHA-256 hex
	PairID       uuid.UUID `gorm:"type:uuid;index:idx_outbox_pair_seq,priority:1"`                       // برای ordering و شاردینگ
	Sequence     uint64    `gorm:"not null;default:0;index:idx_outbox_pair_seq,priority:2"`              // ترتیب قطعی در هر Pair
	Payload      []byte    `gorm:"type:jsonb;not null"`
	Status       string    `gorm:"size:16;not null;default:pending;index:idx_outbox_status_created"`
	CreatedAt    time.Time `gorm:"not null;index:idx_outbox_status_created"`
	SentAt       *time.Time
	ErrorMessage *string `gorm:"size:500"`
	RetryCount   int     `gorm:"not null;default:0"`
}

func (OutboxEvent) TableName() string { return "match_events_outbox" }

func (e *OutboxEvent) UnmarshalPayload(v interface{}) error {
	return json.Unmarshal(e.Payload, v)
}
