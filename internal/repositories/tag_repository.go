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

func (r *TagRepository) CreateTag(tag *models.Tag) error {
	return r.db.Create(tag).Error
}
