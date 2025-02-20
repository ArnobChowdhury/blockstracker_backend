package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID              uint                     `gorm:"primaryKey" json:"id"`
	Name            string                   `gorm:"not null" json:"name"`
	Tasks           []Task                   `gorm:"many2many:task_tags" json:"tasks"`
	RepetitiveTasks []RepetitiveTaskTemplate `gorm:"many2many:repetitive_task_template_tags" json:"repetitiveTasks"`
	CreatedAt       time.Time                `gorm:"autoCreateTime" json:"createdAt"`
	ModifiedAt      time.Time                `gorm:"autoUpdateTime" json:"modifiedAt"`
	DeletedAt       gorm.DeletedAt           `gorm:"index" json:"-"`
}
