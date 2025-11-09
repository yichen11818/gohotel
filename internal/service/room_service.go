package service

import (
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"

	"gorm.io/gorm"
)

// RoomService 房间业务逻辑层
type RoomService struct {
	roomRepo *repository.RoomRepository
}

// NewRoomService 创建房间服务实例
func NewRoomService(roomRepo *repository.RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	RoomNumber    string  `json:"room_number" binding:"required"`
	RoomType      string  `json:"room_type" binding:"required"`
	Floor         int     `json:"floor" binding:"required"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	OriginalPrice float64 `json:"original_price"`
	Capacity      int     `json:"capacity" binding:"required,gt=0"`
	Area          float64 `json:"area"`
	BedType       string  `json:"bed_type"`
	Description   string  `json:"description"`
	Facilities    string  `json:"facilities"`
	Images        string  `json:"images"`
}

// UpdateRoomRequest 更新房间请求
type UpdateRoomRequest struct {
	RoomType      string  `json:"room_type"`
	Floor         int     `json:"floor"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"original_price"`
	Capacity      int     `json:"capacity"`
	Area          float64 `json:"area"`
	BedType       string  `json:"bed_type"`
	Description   string  `json:"description"`
	Facilities    string  `json:"facilities"`
	Images        string  `json:"images"`
	Status        string  `json:"status"`
}

// CreateRoom 创建房间
func (s *RoomService) CreateRoom(req *CreateRoomRequest) (*models.Room, error) {
	// 1. 检查房间号是否已存在
	exists, err := s.roomRepo.ExistsByRoomNumber(req.RoomNumber)
	if err != nil {
		return nil, errors.NewDatabaseError("check room number", err)
	}
	if exists {
		return nil, errors.NewConflictError("房间号已存在")
	}

	// 2. 创建房间对象
	room := &models.Room{
		RoomNumber:    req.RoomNumber,
		RoomType:      req.RoomType,
		Floor:         req.Floor,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Capacity:      req.Capacity,
		Area:          req.Area,
		BedType:       req.BedType,
		Description:   req.Description,
		Facilities:    req.Facilities,
		Images:        req.Images,
		Status:        "available",
	}

	// 3. 保存到数据库
	if err := s.roomRepo.Create(room); err != nil {
		return nil, errors.NewDatabaseError("create room", err)
	}

	return room, nil
}

// GetRoomByID 根据 ID 获取房间
func (s *RoomService) GetRoomByID(id uint) (*models.Room, error) {
	room, err := s.roomRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("房间不存在")
		}
		return nil, errors.NewDatabaseError("find room", err)
	}
	return room, nil
}

// UpdateRoom 更新房间信息
func (s *RoomService) UpdateRoom(id uint, req *UpdateRoomRequest) (*models.Room, error) {
	// 1. 查找房间
	room, err := s.roomRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("房间不存在")
		}
		return nil, errors.NewDatabaseError("find room", err)
	}

	// 2. 更新字段（只更新非空字段）
	if req.RoomType != "" {
		room.RoomType = req.RoomType
	}
	if req.Floor > 0 {
		room.Floor = req.Floor
	}
	if req.Price > 0 {
		room.Price = req.Price
	}
	if req.OriginalPrice > 0 {
		room.OriginalPrice = req.OriginalPrice
	}
	if req.Capacity > 0 {
		room.Capacity = req.Capacity
	}
	if req.Area > 0 {
		room.Area = req.Area
	}
	if req.BedType != "" {
		room.BedType = req.BedType
	}
	if req.Description != "" {
		room.Description = req.Description
	}
	if req.Facilities != "" {
		room.Facilities = req.Facilities
	}
	if req.Images != "" {
		room.Images = req.Images
	}
	if req.Status != "" {
		room.Status = req.Status
	}

	// 3. 保存更新
	if err := s.roomRepo.Update(room); err != nil {
		return nil, errors.NewDatabaseError("update room", err)
	}

	return room, nil
}

// DeleteRoom 删除房间
func (s *RoomService) DeleteRoom(id uint) error {
	// 1. 检查房间是否存在
	_, err := s.roomRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("房间不存在")
		}
		return errors.NewDatabaseError("find room", err)
	}

	// 2. 删除房间
	if err := s.roomRepo.Delete(id); err != nil {
		return errors.NewDatabaseError("delete room", err)
	}

	return nil
}

// ListRooms 获取所有房间列表（分页）
func (s *RoomService) ListRooms(page, pageSize int) ([]models.Room, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	rooms, total, err := s.roomRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("list rooms", err)
	}

	return rooms, total, nil
}

// ListAvailableRooms 获取可用房间列表（分页）
func (s *RoomService) ListAvailableRooms(page, pageSize int) ([]models.Room, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	rooms, total, err := s.roomRepo.FindAvailable(page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("list available rooms", err)
	}

	return rooms, total, nil
}

// SearchRoomsByType 根据房型搜索房间
func (s *RoomService) SearchRoomsByType(roomType string, page, pageSize int) ([]models.Room, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	rooms, total, err := s.roomRepo.FindByRoomType(roomType, page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("search rooms by type", err)
	}

	return rooms, total, nil
}

// SearchRoomsByPrice 根据价格范围搜索房间
func (s *RoomService) SearchRoomsByPrice(minPrice, maxPrice float64, page, pageSize int) ([]models.Room, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	rooms, total, err := s.roomRepo.FindByPriceRange(minPrice, maxPrice, page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("search rooms by price", err)
	}

	return rooms, total, nil
}

// UpdateRoomStatus 更新房间状态
func (s *RoomService) UpdateRoomStatus(id uint, status string) error {
	// 验证状态值
	validStatuses := map[string]bool{
		"available":   true,
		"occupied":    true,
		"maintenance": true,
	}
	if !validStatuses[status] {
		return errors.NewBadRequestError("无效的房间状态")
	}

	if err := s.roomRepo.UpdateStatus(id, status); err != nil {
		return errors.NewDatabaseError("update room status", err)
	}

	return nil
}

