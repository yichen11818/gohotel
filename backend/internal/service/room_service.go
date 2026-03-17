package service

import (
	"context"
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"

	"gorm.io/gorm"
)

// RoomService 房间业务逻辑层
type RoomService struct {
	roomRepo repository.RoomRepository
}

// NewRoomService 创建房间服务实例
func NewRoomService(roomRepo repository.RoomRepository) *RoomService {
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
	Left          int     `json:"left"`
	Top           int     `json:"top"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
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
	Left          int     `json:"left"`
	Top           int     `json:"top"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
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
		Left:          req.Left,
		Top:           req.Top,
		Width:         req.Width,
		Height:        req.Height,
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
	if req.Left > 0 {
		room.Left = req.Left
	}
	if req.Top > 0 {
		room.Top = req.Top
	}
	if req.Width > 0 {
		room.Width = req.Width
	}
	if req.Height > 0 {
		room.Height = req.Height
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

// ListRoomsByFloor 根据楼层获取房间
func (s *RoomService) ListRoomsByFloor(floor, page, pageSize int) ([]models.Room, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	// 直接查询房间表，楼层号相同的一层，无需再查其他表
	rooms, total, err := s.roomRepo.FindRoomByFloor(floor, page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("list rooms by floor", err)
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
func (s *RoomService) UpdateRoomStatus(ctx context.Context, id int64, status string) error {
	// 验证状态值
	validStatuses := map[string]bool{
		"available":   true,
		"occupied":    true,
		"maintenance": true,
	}
	if !validStatuses[status] {
		return errors.NewBadRequestError("无效的房间状态")
	}

	if err := s.roomRepo.UpdateStatus(ctx, id, status); err != nil {
		return errors.NewDatabaseError("update room status", err)
	}

	return nil
}

// UpdateRoomCleanStatus 更新房间清洁状态
func (s *RoomService) UpdateRoomCleanStatus(ctx context.Context, id int64, cleanStatus string) error {
	validStatuses := map[string]bool{
		"clean":      true,
		"dirty":      true,
		"inspecting": true,
	}
	if !validStatuses[cleanStatus] {
		return errors.NewBadRequestError("无效的清洁状态")
	}

	if err := s.roomRepo.UpdateCleanStatus(ctx, id, cleanStatus); err != nil {
		return errors.NewDatabaseError("update room clean status", err)
	}

	return nil
}

// BatchCreateRoomRequest 批量创建房间请求
type BatchCreateRoomRequest struct {
	Rooms []CreateRoomRequest `json:"rooms" binding:"required,min=1,dive"`
}

// BatchCreateRoomsResult 批量创建结果
type BatchCreateRoomsResult struct {
	SuccessCount int            `json:"success_count"` // 成功创建的数量
	FailedCount  int            `json:"failed_count"`  // 失败的数量
	CreatedRooms []*models.Room `json:"created_rooms"` // 成功创建的房间
	FailedRooms  []FailedRoom   `json:"failed_rooms"`  // 失败的房间信息
}

// FailedRoom 创建失败的房间信息
type FailedRoom struct {
	RoomNumber string `json:"room_number"`
	Reason     string `json:"reason"`
}

// BatchCreateRooms 批量创建房间
func (s *RoomService) BatchCreateRooms(req *BatchCreateRoomRequest) (*BatchCreateRoomsResult, error) {
	if len(req.Rooms) == 0 {
		return nil, errors.NewBadRequestError("房间列表不能为空")
	}

	if len(req.Rooms) > 100 {
		return nil, errors.NewBadRequestError("单次最多创建100个房间")
	}

	result := &BatchCreateRoomsResult{
		CreatedRooms: make([]*models.Room, 0),
		FailedRooms:  make([]FailedRoom, 0),
	}

	// 1. 收集所有房间号
	roomNumbers := make([]string, len(req.Rooms))
	for i, r := range req.Rooms {
		roomNumbers[i] = r.RoomNumber
	}

	// 2. 批量检查哪些房间号已存在
	existingNumbers, err := s.roomRepo.ExistsByRoomNumbers(roomNumbers)
	if err != nil {
		return nil, errors.NewDatabaseError("check existing room numbers", err)
	}
	existingSet := make(map[string]bool)
	for _, num := range existingNumbers {
		existingSet[num] = true
	}

	// 3. 检查请求中是否有重复的房间号
	seenNumbers := make(map[string]bool)
	var roomsToCreate []*models.Room

	for _, r := range req.Rooms {
		// 检查请求内是否重复
		if seenNumbers[r.RoomNumber] {
			result.FailedRooms = append(result.FailedRooms, FailedRoom{
				RoomNumber: r.RoomNumber,
				Reason:     "请求中房间号重复",
			})
			result.FailedCount++
			continue
		}
		seenNumbers[r.RoomNumber] = true

		// 检查是否已存在于数据库
		if existingSet[r.RoomNumber] {
			result.FailedRooms = append(result.FailedRooms, FailedRoom{
				RoomNumber: r.RoomNumber,
				Reason:     "房间号已存在",
			})
			result.FailedCount++
			continue
		}

		// 创建房间对象
		room := &models.Room{
			RoomNumber:    r.RoomNumber,
			RoomType:      r.RoomType,
			Floor:         r.Floor,
			Price:         r.Price,
			OriginalPrice: r.OriginalPrice,
			Capacity:      r.Capacity,
			Area:          r.Area,
			BedType:       r.BedType,
			Description:   r.Description,
			Facilities:    r.Facilities,
			Images:        r.Images,
			Status:        "available",
		}
		roomsToCreate = append(roomsToCreate, room)
	}

	// 4. 批量创建有效的房间
	if len(roomsToCreate) > 0 {
		if err := s.roomRepo.CreateBatch(roomsToCreate); err != nil {
			return nil, errors.NewDatabaseError("batch create rooms", err)
		}
		result.CreatedRooms = roomsToCreate
		result.SuccessCount = len(roomsToCreate)
	}

	return result, nil
}
