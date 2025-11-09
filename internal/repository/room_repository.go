package repository

import (
	"gohotel/internal/models"

	"gorm.io/gorm"
)

// RoomRepository 房间数据访问层
type RoomRepository struct {
	db *gorm.DB
}

// NewRoomRepository 创建房间仓库实例
func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// Create 创建房间
func (r *RoomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

// FindByID 根据 ID 查找房间
func (r *RoomRepository) FindByID(id uint) (*models.Room, error) {
	var room models.Room
	err := r.db.First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// FindByRoomNumber 根据房间号查找房间
func (r *RoomRepository) FindByRoomNumber(roomNumber string) (*models.Room, error) {
	var room models.Room
	err := r.db.Where("room_number = ?", roomNumber).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// Update 更新房间信息
func (r *RoomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

// Delete 删除房间
func (r *RoomRepository) Delete(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

// FindAll 查询所有房间（分页）
func (r *RoomRepository) FindAll(page, pageSize int) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	if err := r.db.Model(&models.Room{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("room_number").Find(&rooms).Error
	return rooms, total, err
}

// FindAvailable 查询可用房间（分页）
func (r *RoomRepository) FindAvailable(page, pageSize int) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	query := r.db.Model(&models.Room{}).Where("status = ?", "available")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("price").Find(&rooms).Error
	return rooms, total, err
}

// FindByRoomType 根据房型查询房间（分页）
func (r *RoomRepository) FindByRoomType(roomType string, page, pageSize int) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	query := r.db.Model(&models.Room{}).Where("room_type = ?", roomType)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("price").Find(&rooms).Error
	return rooms, total, err
}

// FindByPriceRange 根据价格范围查询房间（分页）
func (r *RoomRepository) FindByPriceRange(minPrice, maxPrice float64, page, pageSize int) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	query := r.db.Model(&models.Room{}).Where("price BETWEEN ? AND ?", minPrice, maxPrice)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("price").Find(&rooms).Error
	return rooms, total, err
}

// UpdateStatus 更新房间状态
func (r *RoomRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Room{}).Where("id = ?", id).Update("status", status).Error
}

// ExistsByRoomNumber 检查房间号是否已存在
func (r *RoomRepository) ExistsByRoomNumber(roomNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Room{}).Where("room_number = ?", roomNumber).Count(&count).Error
	return count > 0, err
}
