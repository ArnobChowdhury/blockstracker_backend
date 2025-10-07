package repositories

import (
	"blockstracker_backend/models"
	"fmt"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type ChangeRepository struct {
	db *gorm.DB
}

func NewChangeRepository(db *gorm.DB) *ChangeRepository {
	return &ChangeRepository{db: db}
}

// CreateChange creates a new change record within a given transaction,
// calculating the next per-user change_id.
// It requires a transaction object `tx` to ensure atomicity with other database operations.
// NOTE: This implementation assumes the 'changes' table's 'change_id' column is BIGINT NOT NULL
// and part of a composite primary key (user_id, change_id), NOT BIGSERIAL PRIMARY KEY.
func (r *ChangeRepository) CreateChange(tx *gorm.DB, change *models.Change) error {
	// Find the latest change_id for this user within the current transaction
	var latestChangeID int64
	err := tx.Model(&models.Change{}).
		Where("user_id = ?", change.UserID).
		Select("COALESCE(MAX(change_id), 0)"). // COALESCE handles cases where no changes exist for the user yet
		Row().Scan(&latestChangeID)
	if err != nil {
		return fmt.Errorf("failed to get latest change ID for user %s: %w", change.UserID, err)
	}

	// Set the next change ID for the new record
	change.ChangeID = latestChangeID + 1

	return tx.Create(change).Error
}

func (r *ChangeRepository) GetChangesSince(db *gorm.DB, userID uuid.UUID, lastChangeID int64) ([]models.Change, error) {
	var changes []models.Change
	if err := db.Where("user_id = ? AND change_id > ?", userID, lastChangeID).Order("change_id asc").Find(&changes).Error; err != nil {
		return nil, err
	}

	return changes, nil
}
