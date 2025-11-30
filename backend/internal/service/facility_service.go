package service

import (
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"
)

// FacilityService 设施服务层
type FacilityService struct {
	facilityRepo *repository.FacilityRepository
}

// NewFacilityService 创建设施服务实例
func NewFacilityService(facilityRepo *repository.FacilityRepository) *FacilityService {
	return &FacilityService{facilityRepo: facilityRepo}
}

// CreateFacilityRequest 创建设施请求
type CreateFacilityRequest struct {
	Type     string `json:"type" binding:"required"`
	Floor    int    `json:"floor" binding:"required"`
	Left     int    `json:"left" binding:"required"`
	Top      int    `json:"top" binding:"required"`
	Width    int    `json:"width" binding:"required"`
	Height   int    `json:"height" binding:"required"`
	Rotation int    `json:"rotation"`
	Label    string `json:"label"`
}

// UpdateFacilityRequest 更新设施请求
type UpdateFacilityRequest struct {
	Type     string `json:"type"`
	Floor    int    `json:"floor"`
	Left     int    `json:"left"`
	Top      int    `json:"top"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Rotation int    `json:"rotation"`
	Label    string `json:"label"`
}

// BatchUpdateFacilityItem 批量更新设施项
type BatchUpdateFacilityItem struct {
	ID       uint `json:"id" binding:"required"`
	Left     int  `json:"left"`
	Top      int  `json:"top"`
	Width    int  `json:"width"`
	Height   int  `json:"height"`
	Rotation int  `json:"rotation"`
}

// BatchUpdateFacilitiesRequest 批量更新设施请求
type BatchUpdateFacilitiesRequest struct {
	Items []BatchUpdateFacilityItem `json:"items" binding:"required"`
}

// CreateFacility 创建设施
func (s *FacilityService) CreateFacility(req *CreateFacilityRequest) (*models.Facility, error) {
	facility := &models.Facility{
		Type:     req.Type,
		Floor:    req.Floor,
		Left:     req.Left,
		Top:      req.Top,
		Width:    req.Width,
		Height:   req.Height,
		Rotation: req.Rotation,
		Label:    req.Label,
	}
	if err := s.facilityRepo.Create(facility); err != nil {
		return nil, errors.NewDatabaseError("create facility", err)
	}
	return facility, nil
}

// UpdateFacility 更新设施
func (s *FacilityService) UpdateFacility(id uint, req *UpdateFacilityRequest) (*models.Facility, error) {
	facility, err := s.facilityRepo.FindByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("设施不存在")
	}

	// 更新所有字段
	if req.Type != "" {
		facility.Type = req.Type
	}
	if req.Floor != 0 {
		facility.Floor = req.Floor
	}
	if req.Left != 0 {
		facility.Left = req.Left
	}
	if req.Top != 0 {
		facility.Top = req.Top
	}
	if req.Width != 0 {
		facility.Width = req.Width
	}
	if req.Height != 0 {
		facility.Height = req.Height
	}
	facility.Rotation = req.Rotation
	facility.Label = req.Label

	if err := s.facilityRepo.Update(facility); err != nil {
		return nil, errors.NewDatabaseError("update facility", err)
	}
	return facility, nil
}

// DeleteFacility 删除设施
func (s *FacilityService) DeleteFacility(id uint) error {
	// 先检查设施是否存在
	_, err := s.facilityRepo.FindByID(id)
	if err != nil {
		return errors.NewNotFoundError("设施不存在")
	}
	if err := s.facilityRepo.Delete(id); err != nil {
		return errors.NewDatabaseError("delete facility", err)
	}
	return nil
}

// FindFacilityByID 根据 ID 查找设施
func (s *FacilityService) FindFacilityByID(id uint) (*models.Facility, error) {
	facility, err := s.facilityRepo.FindByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError("设施不存在")
	}
	return facility, nil
}

// FindAllFacilities 查询所有设施（分页）
func (s *FacilityService) FindAllFacilities(page, pageSize int) ([]models.Facility, int64, error) {
	facilities, total, err := s.facilityRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("find all facilities", err)
	}
	return facilities, total, nil
}

// FindFacilitiesByFloor 按楼层查询设施
func (s *FacilityService) FindFacilitiesByFloor(floor int) ([]models.Facility, error) {
	facilities, err := s.facilityRepo.FindByFloor(floor)
	if err != nil {
		return nil, errors.NewDatabaseError("find facilities by floor", err)
	}
	return facilities, nil
}

// BatchUpdateFacilities 批量更新设施位置
func (s *FacilityService) BatchUpdateFacilities(req *BatchUpdateFacilitiesRequest) error {
	facilities := make([]models.Facility, len(req.Items))
	for i, item := range req.Items {
		facilities[i] = models.Facility{
			ID:       item.ID,
			Left:     item.Left,
			Top:      item.Top,
			Width:    item.Width,
			Height:   item.Height,
			Rotation: item.Rotation,
		}
	}
	if err := s.facilityRepo.BatchUpdate(facilities); err != nil {
		return errors.NewDatabaseError("batch update facilities", err)
	}
	return nil
}
