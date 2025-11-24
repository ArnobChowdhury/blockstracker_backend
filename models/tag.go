package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string         `json:"name" binding:"required"`
	CreatedAt    JSONTime       `json:"createdAt" binding:"required"`
	ModifiedAt   JSONTime       `json:"modifiedAt" binding:"required"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uuid.UUID      `gorm:"type:uuid;index" json:"userId"`
	LastChangeID int64          `gorm:"not null;default:0" json:"lastChangeId"`
}

type TagRequest struct {
	ID         uuid.UUID `json:"id" binding:"required,uuid"`
	Name       string    `json:"name" binding:"required"`
	CreatedAt  JSONTime  `json:"createdAt" binding:"required"`
	ModifiedAt JSONTime  `json:"modifiedAt" binding:"required"`
}

// Create Tag success response for swagger doc
type TagResponseForSwagger struct {
	Result Tag `json:"result"`
	SuccessResult
}
