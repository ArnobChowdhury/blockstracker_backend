package repositories

import (
	"blockstracker_backend/models"

	"gorm.io/gorm"
)

type SpaceRepository struct {
	db *gorm.DB
}

func NewSpaceRepository(db *gorm.DB) *SpaceRepository {
	return &SpaceRepository{db: db}
}

func (r *SpaceRepository) CreateSpace(Space *models.Space) error {
	return r.db.Create(Space).Error
}
