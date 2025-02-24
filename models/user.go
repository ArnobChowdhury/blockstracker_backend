package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email      string         `gorm:"not null;unique" json:"email"`
	Password   *string        `gorm:"type:varchar" json:"-"`        // Nullable, hidden in JSON
	Provider   *string        `gorm:"type:varchar" json:"provider"` // Nullable
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	ModifiedAt time.Time      `gorm:"autoUpdateTime" json:"modifiedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,strongpassword"`
}
