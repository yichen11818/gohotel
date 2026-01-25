package service

import (
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"

	"gorm.io/gorm"
)

type HotelService struct {
	hotelRepo         *repository.HotelRepository
	hotelSettingsRepo *repository.HotelSettingsRepository
	db                *gorm.DB
}

func NewHotelService(db *gorm.DB, hotelRepo *repository.HotelRepository, hotelSettingsRepo *repository.HotelSettingsRepository) *HotelService {
	return &HotelService{db: db, hotelRepo: hotelRepo, hotelSettingsRepo: hotelSettingsRepo}
}

type CreateHotelRequest struct {
	Name   string  `json:"name" binding:"required"`
	Status *string `json:"status"`
}

type UpdateHotelRequest struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
}

func (s *HotelService) CreateHotel(req *CreateHotelRequest) (*models.Hotel, error) {
	hotel := &models.Hotel{Name: req.Name}
	hotel.Status = "active"
	if req.Status != nil {
		hotel.Status = *req.Status
	}
	if err := s.hotelRepo.Create(hotel); err != nil {
		return nil, errors.NewDatabaseError("create hotel", err)
	}
	return hotel, nil
}

func (s *HotelService) GetHotelByID(id int64) (*models.Hotel, error) {
	hotel, err := s.hotelRepo.FindByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("酒店不存在")
	}
	return hotel, nil
}

func (s *HotelService) ListHotels(page, pageSize int) ([]models.Hotel, int64, error) {
	hotels, total, err := s.hotelRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("list hotels", err)
	}
	return hotels, total, nil
}

func (s *HotelService) UpdateHotel(id int64, req *UpdateHotelRequest) (*models.Hotel, error) {
	hotel, err := s.hotelRepo.FindByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("酒店不存在")
	}
	if req.Name != nil {
		hotel.Name = *req.Name
	}
	if req.Status != nil {
		hotel.Status = *req.Status
	}
	if err := s.hotelRepo.Update(hotel); err != nil {
		return nil, errors.NewDatabaseError("update hotel", err)
	}
	return hotel, nil
}

func (s *HotelService) DeleteHotel(id int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		hsRepo := repository.NewHotelSettingsRepository(tx)
		hRepo := repository.NewHotelRepository(tx)

		_ = hsRepo.DeleteByHotelID(id)
		if err := hRepo.Delete(id); err != nil {
			return errors.NewDatabaseError("delete hotel", err)
		}
		return nil
	})
}
