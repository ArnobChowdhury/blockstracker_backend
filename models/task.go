package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID                       uuid.UUID               `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
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
	ModifiedAt               time.Time               `gorm:"type:timestamp" json:"modifiedAt"`
	Tags                     []Tag                   `json:"tags"`
	Space                    *Space                  `json:"space"`
	SpaceId                  int                     `json:"spaceId"`
	DeletedAt                gorm.DeletedAt          `gorm:"index" json:"-"`
}

type RepetitiveTaskTemplate struct {
	ID                       uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	IsActive                 bool           `gorm:"default:true" json:"isActive"`
	Title                    string         `gorm:"not null" json:"title"`
	Description              *string        `json:"description"`
	Schedule                 string         `gorm:"not null" json:"schedule"`
	Priority                 int            `gorm:"default:3" json:"priority"`
	ShouldBeScored           *bool          `gorm:"default:false" json:"shouldBeScored"`
	Monday                   *bool          `gorm:"default:false" json:"monday"`
	Tuesday                  *bool          `gorm:"default:false" json:"tuesday"`
	Wednesday                *bool          `gorm:"default:false" json:"wednesday"`
	Thursday                 *bool          `gorm:"default:false" json:"thursday"`
	Friday                   *bool          `gorm:"default:false" json:"friday"`
	Saturday                 *bool          `gorm:"default:false" json:"saturday"`
	Sunday                   *bool          `gorm:"default:false" json:"sunday"`
	TimeOfDay                *string        `json:"timeOfDay"`
	LastDateOfTaskGeneration *time.Time     `json:"lastDateOfTaskGeneration"`
	CreatedAt                time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	ModifiedAt               time.Time      `gorm:"type:timestamp" json:"modifiedAt"`
	Tags                     []Tag          `gorm:"many2many:repetitive_task_template_tags" json:"tags"`
	Tasks                    []Task         `json:"tasks"`
	Space                    *Space         `json:"space"`
	SpaceID                  *uint          `json:"spaceId"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
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
