package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string         `json:"name"`
	CreatedAt  time.Time      `json:"createdAt"`
	ModifiedAt time.Time      `json:"modifiedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uuid.UUID      `gorm:"type:uuid;index" json:"userId"`
}

type CreateTagRequest struct {
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}
