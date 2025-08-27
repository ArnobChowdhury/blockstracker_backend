package models

import (
	"time"

	"github.com/google/uuid"
)

type Change struct {
	ChangeID   int64     `gorm:"primaryKey" json:"changeId"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	EntityType string    `gorm:"not null" json:"entityType"`
	EntityID   uuid.UUID `gorm:"type:uuid;not null" json:"entityId"`
	Operation  string    `gorm:"not null" json:"operation"`
	ChangedAt  time.Time `gorm:"not null;default:now()" json:"changedAt"`
}

func (Change) TableName() string {
	return "changes"
}
