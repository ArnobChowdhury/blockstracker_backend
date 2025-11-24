package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRequest struct {
	ID                       uuid.UUID  `json:"id" binding:"required,uuid"`
	IsActive                 *bool      `json:"isActive" binding:"required"`
	Title                    string     `json:"title" binding:"required"`
	Description              string     `json:"description"`
	Schedule                 string     `json:"schedule" binding:"required"`
	Priority                 *int       `json:"priority" binding:"required"`
	CompletionStatus         string     `json:"completionStatus" binding:"required"`
	DueDate                  *JSONTime  `json:"dueDate"`
	ShouldBeScored           *bool      `json:"shouldBeScored" binding:"required"`
	Score                    *int       `json:"score"`
	TimeOfDay                *string    `json:"timeOfDay"`
	RepetitiveTaskTemplateID *uuid.UUID `json:"repetitiveTaskTemplateId"`
	CreatedAt                JSONTime   `json:"createdAt" binding:"required"`
	ModifiedAt               JSONTime   `json:"modifiedAt" binding:"required"`
	Tags                     []Tag      `gorm:"many2many:task_tags;" json:"tags"`
	SpaceID                  *uuid.UUID `gorm:"type:uuid" json:"spaceId"`
}

type Task struct {
	ID                       uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	IsActive                 bool       `gorm:"default:true" json:"isActive"`
	Title                    string     `gorm:"not null" json:"title"`
	Description              string     `json:"description"`
	Schedule                 string     `json:"schedule"`
	Priority                 int        `gorm:"default:3" json:"priority"`
	CompletionStatus         string     `gorm:"default:'INCOMPLETE'" json:"completionStatus"`
	DueDate                  *JSONTime  `json:"dueDate"`
	ShouldBeScored           *bool      `json:"shouldBeScored"`
	Score                    *int       `json:"score"`
	TimeOfDay                *string    `json:"timeOfDay"`
	RepetitiveTaskTemplateID *uuid.UUID `gorm:"type:uuid" json:"repetitiveTaskTemplateId"`
	CreatedAt                JSONTime   `json:"createdAt"`
	ModifiedAt               JSONTime   `json:"modifiedAt"`
	// Tags                     []Tag          `gorm:"many2many:task_tags;" json:"tags"`
	SpaceID      *uuid.UUID     `gorm:"type:uuid" json:"spaceId"`
	UserID       uuid.UUID      `gorm:"type:uuid" json:"userId"` // Add UserID here
	LastChangeID int64          `gorm:"not null;default:0" json:"lastChangeId"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Create Task success response for swagger doc
type TaskResponseForSwagger struct {
	Result Task `json:"result"`
	SuccessResult
}

type RepetitiveTaskTemplate struct {
	ID                       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	IsActive                 bool      `gorm:"default:true" json:"isActive"`
	Title                    string    `gorm:"not null" json:"title"`
	Description              *string   `json:"description"`
	Schedule                 string    `gorm:"not null" json:"schedule"`
	Priority                 int       `gorm:"default:3" json:"priority"`
	ShouldBeScored           *bool     `gorm:"default:false" json:"shouldBeScored"`
	Monday                   *bool     `gorm:"default:false" json:"monday"`
	Tuesday                  *bool     `gorm:"default:false" json:"tuesday"`
	Wednesday                *bool     `gorm:"default:false" json:"wednesday"`
	Thursday                 *bool     `gorm:"default:false" json:"thursday"`
	Friday                   *bool     `gorm:"default:false" json:"friday"`
	Saturday                 *bool     `gorm:"default:false" json:"saturday"`
	Sunday                   *bool     `gorm:"default:false" json:"sunday"`
	TimeOfDay                *string   `json:"timeOfDay"`
	LastDateOfTaskGeneration *JSONTime `json:"lastDateOfTaskGeneration"`
	CreatedAt                JSONTime  `json:"createdAt"`
	ModifiedAt               JSONTime  `json:"modifiedAt"`
	// Tags                     []Tag          `gorm:"many2many:repetitive_task_template_tags" json:"tags"`
	SpaceID      *uuid.UUID     `gorm:"type:uuid" json:"spaceId"`
	UserID       uuid.UUID      `gorm:"type:uuid" json:"userId"` // Add UserID here
	LastChangeID int64          `gorm:"not null;default:0" json:"lastChangeId"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type RepetitiveTaskTemplateRequest struct {
	ID                       uuid.UUID      `json:"id" binding:"required,uuid"`
	IsActive                 *bool          `json:"isActive" binding:"required"`
	Title                    string         `json:"title" binding:"required"`
	Description              *string        `json:"description"`
	Schedule                 string         `json:"schedule" binding:"required"`
	Priority                 *int           `json:"priority" binding:"required"`
	ShouldBeScored           *bool          `json:"shouldBeScored" binding:"required"`
	Monday                   *bool          `json:"monday" binding:"required"`
	Tuesday                  *bool          `json:"tuesday" binding:"required"`
	Wednesday                *bool          `json:"wednesday" binding:"required"`
	Thursday                 *bool          `json:"thursday" binding:"required"`
	Friday                   *bool          `json:"friday" binding:"required"`
	Saturday                 *bool          `json:"saturday" binding:"required"`
	Sunday                   *bool          `json:"sunday" binding:"required"`
	TimeOfDay                *string        `json:"timeOfDay"`
	LastDateOfTaskGeneration *JSONTime      `json:"lastDateOfTaskGeneration"`
	CreatedAt                JSONTime       `json:"createdAt" binding:"required"`
	ModifiedAt               JSONTime       `json:"modifiedAt" binding:"required"`
	Tags                     []Tag          `json:"tags"`
	Tasks                    []Task         `json:"tasks"`
	SpaceID                  *uuid.UUID     `json:"spaceId"`
	DeletedAt                gorm.DeletedAt `json:"-"` // Keep json:"-" to omit from JSON
}

type RepetitiveTaskTemplateResponseForSwagger struct {
	Result RepetitiveTaskTemplate `json:"result"`
	SuccessResult
}

type UpdateRepetitiveTaskRequest struct {
	TaskID    int    `json:"task_id" binding:"required"`
	Frequency string `json:"frequency" binding:"required"`
}
