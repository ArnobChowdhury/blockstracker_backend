package repositories

import (
	"blockstracker_backend/models"

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
