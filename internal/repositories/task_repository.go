package repositories

import (
	"blockstracker_backend/models"
	"time"

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
	if err := tx.Model(&models.Task{}).Preload("Tags").Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetTaskByRepetitiveTemplateIDAndDueDate(tx *gorm.DB, templateID uuid.UUID, dueDate time.Time, userID uuid.UUID) (*models.Task, error) {
	var task models.Task
	if err := tx.Model(&models.Task{}).Where("repetitive_task_template_id = ? AND due_date = ? AND user_id = ?", templateID, dueDate, userID).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetTasksByIDs(tx *gorm.DB, taskIDs []uuid.UUID, userID uuid.UUID) ([]models.Task, error) {
	var tasks []models.Task
	if err := tx.Model(&models.Task{}).Preload("Tags").Where("id IN ? AND user_id = ?", taskIDs, userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
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
	if err := tx.Model(&models.RepetitiveTaskTemplate{}).Preload("Tags").Where("id = ? AND user_id = ?", templateID, userID).First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *TaskRepository) GetRepetitiveTaskTemplatesByIDs(tx *gorm.DB, templateIDs []uuid.UUID, userID uuid.UUID) ([]models.RepetitiveTaskTemplate, error) {
	var templates []models.RepetitiveTaskTemplate
	if err := tx.Model(&models.RepetitiveTaskTemplate{}).Preload("Tags").Where("id IN ? AND user_id = ?", templateIDs, userID).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *TaskRepository) UpdateRepetitiveTaskTemplate(tx *gorm.DB, templateID, userID uuid.UUID, data map[string]any) error {
	result := tx.Model(&models.RepetitiveTaskTemplate{}).Where("id = ? AND user_id = ?", templateID, userID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
