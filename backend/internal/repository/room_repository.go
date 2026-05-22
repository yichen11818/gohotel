package repository

import (
	"context"
	"gohotel/internal/models"
	"strings"

	"gorm.io/gorm"
)

// RoomRepository 房间数据访问接口
type RoomRepository interface {
	Create(room *models.Room) error
	FindByID(id uint) (*models.Room, error)
	FindByRoomNumber(roomNumber string) (*models.Room, error)
	Update(room *models.Room) error
	Delete(id uint) error
	FindAll(page, pageSize int, roomNumber, roomType, status, cleanStatus string) ([]models.Room, int64, error)
	FindAvailable(page, pageSize int) ([]models.Room, int64, error)
	FindByRoomType(roomType string, page, pageSize int) ([]models.Room, int64, error)
	FindByPriceRange(minPrice, maxPrice float64, page, pageSize int) ([]models.Room, int64, error)
	FindRoomByFloor(floor, page, pageSize int) ([]models.Room, int64, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	UpdateCleanStatus(ctx context.Context, id int64, cleanStatus string) error
	ExistsByRoomNumber(roomNumber string) (bool, error)
	CreateBatch(rooms []*models.Room) error
	ExistsByRoomNumbers(roomNumbers []string) ([]string, error)
}

type roomRepository struct {
	db *gorm.DB
}

// NewRoomRepository 创建房间仓库实例
func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

// Create 创建房间
func (r *roomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

// FindByID 根据 ID 查找房间
func (r *roomRepository) FindByID(id uint) (*models.Room, error) {
	var room models.Room
	err := r.db.First(&room, id).Error
	if err != nil {
		return nil, err
	}
	if err := r.hydrateRooms([]*models.Room{&room}); err != nil {
		return nil, err
	}
	return &room, nil
}

// FindByRoomNumber 根据房间号查找房间
func (r *roomRepository) FindByRoomNumber(roomNumber string) (*models.Room, error) {
	var room models.Room
	err := r.db.Where("room_number = ?", roomNumber).First(&room).Error
	if err != nil {
		return nil, err
	}
	if err := r.hydrateRooms([]*models.Room{&room}); err != nil {
		return nil, err
	}
	return &room, nil
}

// Update 更新房间信息
func (r *roomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

// Delete 删除房间
func (r *roomRepository) Delete(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

// FindAll 查询所有房间（分页）
func (r *roomRepository) FindAll(page, pageSize int, roomNumber, roomType, status, cleanStatus string) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	query := r.db.Model(&models.Room{})
	if roomNumber != "" {
		query = query.Where("room_number LIKE ?", "%"+roomNumber+"%")
	}
	if roomType != "" {
		query = query.Where("room_type = ?", roomType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if cleanStatus != "" {
		query = query.Where("clean_status = ?", cleanStatus)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("room_number").Find(&rooms).Error
	if err == nil {
		err = r.hydrateRoomSlice(rooms)
	}
	return rooms, total, err
}

// FindAvailable 查询可用房间（分页）
func (r *roomRepository) FindAvailable(page, pageSize int) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	query := r.db.Model(&models.Room{}).Where("status = ?", "available")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("price").Find(&rooms).Error
	if err == nil {
		err = r.hydrateRoomSlice(rooms)
	}
	return rooms, total, err
}

// FindByRoomType 根据房型查询房间（分页）
func (r *roomRepository) FindByRoomType(roomType string, page, pageSize int) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	query := r.db.Model(&models.Room{}).Where("room_type = ?", roomType)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("price").Find(&rooms).Error
	if err == nil {
		err = r.hydrateRoomSlice(rooms)
	}
	return rooms, total, err
}

// FindByPriceRange 根据价格范围查询房间（分页）
func (r *roomRepository) FindByPriceRange(minPrice, maxPrice float64, page, pageSize int) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	query := r.db.Model(&models.Room{}).Where("price BETWEEN ? AND ?", minPrice, maxPrice)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("price").Find(&rooms).Error
	if err == nil {
		err = r.hydrateRoomSlice(rooms)
	}
	return rooms, total, err
}

// FindRoomByFloor 根据楼层查询房间
func (r *roomRepository) FindRoomByFloor(floor, page, pageSize int) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	query := r.db.Model(&models.Room{}).Where("floor = ?", floor)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("room_number").Find(&rooms).Error //排序方式
	if err == nil {
		err = r.hydrateRoomSlice(rooms)
	}
	return rooms, total, err
}

// UpdateStatus 更新房间状态
func (r *roomRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&models.Room{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateCleanStatus 更新房间清洁状态
func (r *roomRepository) UpdateCleanStatus(ctx context.Context, id int64, cleanStatus string) error {
	return r.db.WithContext(ctx).Model(&models.Room{}).Where("id = ?", id).Update("clean_status", cleanStatus).Error
}

// ExistsByRoomNumber 检查房间号是否已存在
func (r *roomRepository) ExistsByRoomNumber(roomNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Room{}).Where("room_number = ?", roomNumber).Count(&count).Error
	return count > 0, err
}

// CreateBatch 批量创建房间
func (r *roomRepository) CreateBatch(rooms []*models.Room) error {
	return r.db.Create(rooms).Error
}

// ExistsByRoomNumbers 批量检查房间号是否已存在，返回已存在的房间号列表
func (r *roomRepository) ExistsByRoomNumbers(roomNumbers []string) ([]string, error) {
	var existingRooms []models.Room
	err := r.db.Model(&models.Room{}).Where("room_number IN ?", roomNumbers).Select("room_number").Find(&existingRooms).Error
	if err != nil {
		return nil, err
	}
	var existingNumbers []string
	for _, room := range existingRooms {
		existingNumbers = append(existingNumbers, room.RoomNumber)
	}
	return existingNumbers, nil
}

func (r *roomRepository) hydrateRoomSlice(rooms []models.Room) error {
	if len(rooms) == 0 {
		return nil
	}

	roomRefs := make([]*models.Room, 0, len(rooms))
	for i := range rooms {
		roomRefs = append(roomRefs, &rooms[i])
	}

	return r.hydrateRooms(roomRefs)
}

func (r *roomRepository) hydrateRooms(rooms []*models.Room) error {
	if len(rooms) == 0 {
		return nil
	}

	roomTypes := make([]string, 0, len(rooms))
	seen := make(map[string]struct{}, len(rooms))
	for _, room := range rooms {
		roomType := strings.TrimSpace(room.RoomType)
		if roomType == "" {
			continue
		}
		if _, exists := seen[roomType]; exists {
			continue
		}
		seen[roomType] = struct{}{}
		roomTypes = append(roomTypes, roomType)
	}

	if len(roomTypes) == 0 {
		return nil
	}

	var categories []models.RoomCategory
	if err := r.db.Where("name IN ?", roomTypes).Find(&categories).Error; err != nil {
		return err
	}

	categoryMap := make(map[string]models.RoomCategory, len(categories))
	for _, category := range categories {
		categoryMap[category.Name] = category
	}

	for _, room := range rooms {
		category, exists := categoryMap[room.RoomType]
		if !exists {
			continue
		}
		if category.Description != "" {
			room.Description = category.Description
		}
		if category.Facilities != "" {
			room.Facilities = category.Facilities
		}
		if category.Images != "" {
			room.Images = category.Images
		}
	}

	return nil
}
