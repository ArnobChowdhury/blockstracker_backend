package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Space struct {
	ID              uuid.UUID                `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name            string                   `gorm:"unique;not null" json:"name"`
	Tasks           []Task                   `json:"tasks"`
	RepetitiveTasks []RepetitiveTaskTemplate `json:"repetitiveTasks"`
	CreatedAt       time.Time                `gorm:"autoCreateTime" json:"createdAt"`
	ModifiedAt      time.Time                `gorm:"autoUpdateTime" json:"modifiedAt"`
	DeletedAt       gorm.DeletedAt           `gorm:"index" json:"-"`
}
