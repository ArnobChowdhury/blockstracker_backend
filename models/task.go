package models

import (
	"time"

	"gorm.io/gorm"
)

// Task model definition
type Task struct {
	ID                       uint                    `gorm:"primaryKey" json:"id"`
	IsActive                 bool                    `gorm:"default:true" json:"isActive"`
	Title                    string                  `gorm:"not null" json:"title"`
	Description              string                  `json:"description"`
	Schedule                 string                  `json:"schedule"`
	Priority                 int                     `gorm:"default:3" json:"priority"`
	CompletionStatus         string                  `gorm:"default:'INCOMPLETE'" json:"completionStatus"`
	DueDate                  *time.Time              `json:"dueDate"`
	ShouldBeScored           *bool                   `json:"shouldBeScored"`
	Score                    *int                    `json:"score"`
	TimeOfDay                *string                 `json:"timeOfDay"`
	RepetitiveTaskTemplate   *RepetitiveTaskTemplate `json:"repetitiveTaskTemplate"`
	RepetitiveTaskTemplateId int                     `json:"repetitiveTaskTemplateId"`
	CreatedAt                time.Time               `json:"createdAt"`
	ModifiedAt               time.Time               `json:"modifiedAt"`
	Tags                     []Tag                   `json:"tags"`
	Space                    *Space                  `json:"space"`
	SpaceId                  int                     `json:"spaceId"`

	// Unique constraint across repetitiveTaskTemplateId and dueDate
	// GORM does not enforce this automatically in the database, but you can add it in migrations manually
	gorm.Model
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type UpdateTaskRequest struct {
	TaskID      int    `json:"task_id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type UpdateRepetitiveTaskRequest struct {
	TaskID    int    `json:"task_id" binding:"required"`
	Frequency string `json:"frequency" binding:"required"`
}
