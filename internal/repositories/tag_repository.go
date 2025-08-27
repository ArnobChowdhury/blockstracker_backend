package repositories

import (
	"blockstracker_backend/models"

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
