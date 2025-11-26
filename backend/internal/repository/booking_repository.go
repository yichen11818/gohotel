package repository

import (
	"gohotel/internal/models"
	"time"

	"gorm.io/gorm"
)

// BookingRepository 预订数据访问层
type BookingRepository struct {
	db *gorm.DB
}

// NewBookingRepository 创建预订仓库实例
func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// Create 创建预订
func (r *BookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

// FindByID 根据 ID 查找预订（包含关联的用户和房间信息）
func (r *BookingRepository) FindByID(id int64) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("User").Preload("Room").First(&booking, id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByBookingNumber 根据订单号查找预订
func (r *BookingRepository) FindByBookingNumber(bookingNumber string) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("User").Preload("Room").
		Where("booking_number = ?", bookingNumber).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByUserID 根据用户 ID 查找预订列表（分页）
func (r *BookingRepository) FindByUserID(userID int64, page, pageSize int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Room").Offset(offset).Limit(pageSize).
		Order("created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

// FindByRoomID 根据房间 ID 查找预订列表
func (r *BookingRepository) FindByRoomID(roomID int64) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("room_id = ?", roomID).
		Where("status IN ?", []string{"pending", "confirmed", "checkin"}).
		Order("check_in").Find(&bookings).Error
	return bookings, err
}

// FindByRoomNumberAndStatus 根据房间号和状态查找预订列表
func (r *BookingRepository) FindByRoomNumberAndStatus(roomNumber string, status string) ([]models.Booking, error) {
	var bookings []models.Booking
	query := r.db.Model(&models.Booking{}).
		Joins("JOIN rooms ON rooms.id = bookings.room_id").
		Where("rooms.room_number = ?", roomNumber)

	// 根据状态参数过滤
	if status != "" {
		query = query.Where("bookings.status = ?", status)
	}

	err := query.Preload("User").Preload("Room").
		Order("bookings.created_at DESC").Find(&bookings).Error

	return bookings, err
}

// FindAll 查询所有预订（分页）
func (r *BookingRepository) FindAll(page, pageSize int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	if err := r.db.Model(&models.Booking{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Preload("User").Preload("Room").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

// Update 更新预订信息
func (r *BookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

// UpdateStatus 更新预订状态
func (r *BookingRepository) UpdateStatus(id int64, status string) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).Update("status", status).Error
}

// UpdatePaymentStatus 更新支付状态
func (r *BookingRepository) UpdatePaymentStatus(id int64, paymentStatus string) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).Update("payment_status", paymentStatus).Error
}

// Delete 删除预订
func (r *BookingRepository) Delete(id int64) error {
	return r.db.Delete(&models.Booking{}, id).Error
}

// CheckRoomAvailability 检查房间在指定日期范围内是否可用
// 如果有冲突的预订，返回 false
func (r *BookingRepository) CheckRoomAvailability(roomID int64, checkIn, checkOut time.Time) (bool, error) {
	var count int64

	// 查询是否有冲突的预订
	// 冲突条件：
	// 1. 同一房间
	// 2. 预订未取消
	// 3. 日期有重叠
	err := r.db.Model(&models.Booking{}).
		Where("room_id = ?", roomID).
		Where("status IN ?", []string{"pending", "confirmed", "checkin"}).
		Where("(check_in < ? AND check_out > ?)", checkOut, checkIn).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// FindByDateRange 根据日期范围查询预订（分页）
func (r *BookingRepository) FindByDateRange(startDate, endDate time.Time, page, pageSize int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).
		Where("check_in >= ? AND check_out <= ?", startDate, endDate)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Preload("Room").
		Offset(offset).Limit(pageSize).
		Order("check_in").Find(&bookings).Error
	return bookings, total, err
}

// FindByStatus 根据状态查询预订（分页）
func (r *BookingRepository) FindByStatus(status string, page, pageSize int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Preload("Room").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

// FindByGuestInfo 通过客人姓名、手机号和状态查询预订
func (r *BookingRepository) FindByGuestInfo(guestName, guestPhone, status string) ([]models.Booking, error) {
	var bookings []models.Booking
	query := r.db.Model(&models.Booking{})

	if guestName != "" {
		query = query.Where("guest_name LIKE ?", "%"+guestName+"%")
	}

	if guestPhone != "" {
		query = query.Where("guest_phone LIKE ?", "%"+guestPhone+"%")
	}

	// 根据状态参数过滤，为空则不过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Preload("User").Preload("Room").
		Order("created_at DESC").Find(&bookings).Error
	if err != nil {
		return nil, err
	}

	return bookings, nil
}
