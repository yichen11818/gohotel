package repository

import (
	"gohotel/internal/models"

	"gorm.io/gorm"
)

type HotelRepository struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) *HotelRepository {
	return &HotelRepository{db: db}
}

func (r *HotelRepository) Create(hotel *models.Hotel) error {
	return r.db.Create(hotel).Error
}

func (r *HotelRepository) FindByID(id int64) (*models.Hotel, error) {
	var hotel models.Hotel
	if err := r.db.First(&hotel, id).Error; err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (r *HotelRepository) FindAll(page, pageSize int) ([]models.Hotel, int64, error) {
	var hotels []models.Hotel
	var total int64

	query := r.db.Model(&models.Hotel{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&hotels).Error; err != nil {
		return nil, 0, err
	}

	return hotels, total, nil
}

func (r *HotelRepository) Update(hotel *models.Hotel) error {
	return r.db.Save(hotel).Error
}

func (r *HotelRepository) Delete(id int64) error {
	return r.db.Delete(&models.Hotel{}, id).Error
}
