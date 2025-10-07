package repositories

import (
	"blockstracker_backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SpaceRepository struct {
	db *gorm.DB
}

func NewSpaceRepository(db *gorm.DB) *SpaceRepository {
	return &SpaceRepository{db: db}
}

func (r *SpaceRepository) CreateSpace(tx *gorm.DB, Space *models.Space) error {
	return tx.Create(Space).Error
}

func (r *SpaceRepository) GetSpaceByID(tx *gorm.DB, spaceID uuid.UUID, userID uuid.UUID) (*models.Space, error) {
	var space models.Space
	if err := tx.Model(&models.Space{}).Where("id = ? AND user_id = ?", spaceID, userID).First(&space).Error; err != nil {
		return nil, err
	}
	return &space, nil
}

func (r *SpaceRepository) GetSpacesByIDs(tx *gorm.DB, spaceIDs []uuid.UUID, userID uuid.UUID) ([]models.Space, error) {
	var spaces []models.Space
	if err := tx.Model(&models.Space{}).Where("id IN ? AND user_id = ?", spaceIDs, userID).Find(&spaces).Error; err != nil {
		return nil, err
	}
	return spaces, nil
}

func (r *SpaceRepository) UpdateSpace(tx *gorm.DB, space *models.Space) error {
	result := tx.Model(&models.Space{}).Where("id = ? AND user_id = ?", space.ID, space.UserID).Updates(space)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
