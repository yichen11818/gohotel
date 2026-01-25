package repository

import (
	"gohotel/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HotelSettingsRepository struct {
	db *gorm.DB
}

func NewHotelSettingsRepository(db *gorm.DB) *HotelSettingsRepository {
	return &HotelSettingsRepository{db: db}
}

func (r *HotelSettingsRepository) FindByHotelID(hotelID int64) (*models.HotelSettings, error) {
	var hs models.HotelSettings
	if err := r.db.Where("hotel_id = ?", hotelID).First(&hs).Error; err != nil {
		return nil, err
	}
	return &hs, nil
}

func (r *HotelSettingsRepository) Upsert(hs *models.HotelSettings) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "hotel_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"settings_json", "updated_at"}),
	}).Create(hs).Error
}

func (r *HotelSettingsRepository) DeleteByHotelID(hotelID int64) error {
	return r.db.Where("hotel_id = ?", hotelID).Delete(&models.HotelSettings{}).Error
}
