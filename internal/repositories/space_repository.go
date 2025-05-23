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

func (r *SpaceRepository) UpdateSpace(space *models.Space) error {
	result := r.db.Model(&models.Space{}).Where("id = ? AND user_id = ?", space.ID, space.UserID).Updates(space)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
