package service

import (
	"encoding/json"
	stderrors "errors"
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"

	"gorm.io/gorm"
)

type HotelSettingsService struct {
	db        *gorm.DB
	hsRepo    *repository.HotelSettingsRepository
	hotelRepo *repository.HotelRepository
}

func NewHotelSettingsService(db *gorm.DB, hsRepo *repository.HotelSettingsRepository, hotelRepo *repository.HotelRepository) *HotelSettingsService {
	return &HotelSettingsService{db: db, hsRepo: hsRepo, hotelRepo: hotelRepo}
}

type SaveHotelSettingsRequest struct {
	HotelID  int64                  `json:"hotel_id" binding:"required"`
	Settings map[string]interface{} `json:"settings" binding:"required"`
}

func (s *HotelSettingsService) GetSettingsByHotelID(hotelID int64) (map[string]interface{}, error) {
	_, err := s.hotelRepo.FindByID(hotelID)
	if err != nil {
		return nil, errors.NewNotFoundError("酒店不存在")
	}

	hs, err := s.hsRepo.FindByHotelID(hotelID)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{}, nil
		}
		return nil, errors.NewDatabaseError("get hotel settings", err)
	}

	if hs.SettingsJSON == "" {
		return map[string]interface{}{}, nil
	}

	var out map[string]interface{}
	if err := json.Unmarshal([]byte(hs.SettingsJSON), &out); err != nil {
		return nil, errors.NewInternalServerError("酒店设置解析失败")
	}
	if out == nil {
		out = map[string]interface{}{}
	}
	return out, nil
}

func (s *HotelSettingsService) Save(req *SaveHotelSettingsRequest) error {
	_, err := s.hotelRepo.FindByID(req.HotelID)
	if err != nil {
		return errors.NewNotFoundError("酒店不存在")
	}

	b, err := json.Marshal(req.Settings)
	if err != nil {
		return errors.NewBadRequestError("settings 格式错误")
	}

	hs := &models.HotelSettings{
		HotelID:      req.HotelID,
		SettingsJSON: string(b),
	}

	if err := s.hsRepo.Upsert(hs); err != nil {
		return errors.NewDatabaseError("save hotel settings", err)
	}
	return nil
}

func (s *HotelSettingsService) GetPublicSettingsByHotelID(hotelID int64) (map[string]interface{}, error) {
	return s.GetSettingsByHotelID(hotelID)
}
