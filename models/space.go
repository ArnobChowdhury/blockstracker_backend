package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Space struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string         `json:"name" binding:"required"`
	CreatedAt  time.Time      `json:"createdAt" binding:"required"`
	ModifiedAt time.Time      `json:"modifiedAt" binding:"required"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uuid.UUID      `gorm:"type:uuid;index" json:"userId"`
}

type SpaceRequest struct {
	Name       string    `json:"name" binding:"required"`
	CreatedAt  time.Time `json:"createdAt" binding:"required"`
	ModifiedAt time.Time `json:"modifiedAt" binding:"required"`
}

// Create Space success response for swagger doc
type SpaceResponseForSwagger struct {
	Result Space `json:"result"`
	SuccessResult
}
