package repositories

import (
	"blockstracker_backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(tx *gorm.DB, task *models.Task) error {
	return tx.Create(task).Error
}

func (r *TaskRepository) GetTaskByID(tx *gorm.DB, taskID uuid.UUID, userID uuid.UUID) (*models.Task, error) {
	var task models.Task
	if err := tx.Model(&models.Task{}).Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) UpdateTask(tx *gorm.DB, task *models.Task) error {
	result := tx.Model(&models.Task{}).Where("id = ? AND user_id = ?", task.ID, task.UserID).Updates(task)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TaskRepository) CreateRepetitiveTaskTemplate(tx *gorm.DB, repetitiveTaskTemplate *models.RepetitiveTaskTemplate) error {
	return tx.Create(repetitiveTaskTemplate).Error
}

func (r *TaskRepository) GetRepetitiveTaskTemplateByID(tx *gorm.DB, templateID uuid.UUID, userID uuid.UUID) (*models.RepetitiveTaskTemplate, error) {
	var template models.RepetitiveTaskTemplate
	if err := tx.Model(&models.RepetitiveTaskTemplate{}).Where("id = ? AND user_id = ?", templateID, userID).First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *TaskRepository) UpdateRepetitiveTaskTemplate(tx *gorm.DB, repetitiveTaskTemplate *models.RepetitiveTaskTemplate) error {
	result := tx.Model(&models.RepetitiveTaskTemplate{}).Where("id = ? AND user_id = ?", repetitiveTaskTemplate.ID, repetitiveTaskTemplate.UserID).Updates(repetitiveTaskTemplate)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
