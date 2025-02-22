package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Username   string         `gorm:"not null;unique" json:"username"`
	Email      string         `gorm:"not null;unique" json:"email"`
	Password   string         `gorm:"not null" json:"-"` // Excluded from JSON responses
	Provider   string         `json:"provider"`
	CreatedAt  time.Time      `json:"createdAt"`
	ModifiedAt time.Time      `gorm:"type:timestamp" json:"modifiedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
