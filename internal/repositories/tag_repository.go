package repositories

import (
	"blockstracker_backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) CreateTag(tx *gorm.DB, tag *models.Tag) error {
	return tx.Create(tag).Error
}

func (r *TagRepository) GetTagByID(tx *gorm.DB, tagID uuid.UUID, userID uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	if err := tx.Model(&models.Tag{}).Preload("Tasks").Preload("RepetitiveTasks").Where("id = ? AND user_id = ?", tagID, userID).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) GetTagsByIDs(tx *gorm.DB, tagIDs []uuid.UUID, userID uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag
	if err := tx.Model(&models.Tag{}).Preload("Tasks").Preload("RepetitiveTasks").Where("id IN ? AND user_id = ?", tagIDs, userID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepository) UpdateTag(tx *gorm.DB, tag *models.Tag) error {
	result := tx.Model(&models.Tag{}).Where(
		"id = ? AND user_id = ?", tag.ID, tag.UserID).Updates(tag)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
