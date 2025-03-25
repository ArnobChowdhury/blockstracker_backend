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

func (r *TaskRepository) CreateTask(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) UpdateTask(task *models.Task) error {
	result := r.db.Model(&models.Task{}).Where("id = ? AND user_id = ?", task.ID, task.UserID).Updates(task)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TaskRepository) CreateRepetitiveTaskTemplate(repetitiveTaskTemplate *models.RepetitiveTaskTemplate) error {
	return r.db.Create(repetitiveTaskTemplate).Error
}
