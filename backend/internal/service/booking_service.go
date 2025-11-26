package service

import (
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// BookingService 预订业务逻辑层
type BookingService struct {
	bookingRepo *repository.BookingRepository
	roomRepo    *repository.RoomRepository
	userRepo    *repository.UserRepository
}

// NewBookingService 创建预订服务实例
func NewBookingService(
	bookingRepo *repository.BookingRepository,
	roomRepo *repository.RoomRepository,
	userRepo *repository.UserRepository,
) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
		userRepo:    userRepo,
	}
}

// CreateBookingRequest 创建预订请求
type CreateBookingRequest struct {
	RoomID         int64  `json:"room_id" binding:"required"`
	CheckIn        string `json:"check_in" binding:"required"`  // 格式: "2024-01-01"
	CheckOut       string `json:"check_out" binding:"required"` // 格式: "2024-01-05"
	GuestName      string `json:"guest_name" binding:"required"`
	GuestPhone     string `json:"guest_phone" binding:"required"`
	GuestIDCard    string `json:"guest_id_card"`   // 入住人身份证号，可选
	SpecialRequest string `json:"special_request"` // 特殊要求，可选
}

// CreateBooking 创建预订
func (s *BookingService) CreateBooking(userID int64, req *CreateBookingRequest) (*models.Booking, error) {
	// 1. 验证日期格式
	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		return nil, errors.NewBadRequestError("入住日期格式错误，应为: YYYY-MM-DD")
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		return nil, errors.NewBadRequestError("退房日期格式错误，应为: YYYY-MM-DD")
	}

	// 2. 验证日期逻辑
	now := time.Now().Truncate(24 * time.Hour) // 去掉时间部分，只保留日期
	if checkIn.Before(now) {
		return nil, errors.NewBadRequestError("入住日期不能早于今天")
	}
	if checkOut.Before(checkIn) || checkOut.Equal(checkIn) {
		return nil, errors.NewBadRequestError("退房日期必须晚于入住日期")
	}

	// 3. 查询房间是否存在
	room, err := s.roomRepo.FindByID(uint(req.RoomID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("房间不存在")
		}
		return nil, errors.NewDatabaseError("find room", err)
	}

	// 4. 检查房间状态
	if !room.IsAvailable() {
		return nil, errors.NewBadRequestError("房间不可用")
	}

	// 5. 检查房间在指定日期是否已被预订
	available, err := s.bookingRepo.CheckRoomAvailability(req.RoomID, checkIn, checkOut)
	if err != nil {
		return nil, errors.NewDatabaseError("check availability", err)
	}
	if !available {
		return nil, errors.NewConflictError("该房间在所选日期已被预订")
	}

	// 6. 计算总天数和总价
	totalDays := int(checkOut.Sub(checkIn).Hours() / 24)
	totalPrice := float64(totalDays) * room.Price

	// 7. 生成订单号和预订ID
	bookingNumber := utils.GenID()
	bookingID := utils.GenID()

	// 8. 创建预订对象
	booking := &models.Booking{
		ID:             utils.JSONInt64(bookingID),
		BookingNumber:  utils.JSONInt64(bookingNumber),
		UserID:         utils.JSONInt64(userID),
		RoomID:         req.RoomID,
		CheckIn:        checkIn,
		CheckOut:       checkOut,
		TotalDays:      totalDays,
		TotalPrice:     totalPrice,
		GuestName:      req.GuestName,
		GuestPhone:     req.GuestPhone,
		GuestIDCard:    req.GuestIDCard,
		SpecialRequest: req.SpecialRequest,
		Status:         "pending",
		PaymentStatus:  "unpaid",
	}

	// 9. 保存到数据库
	if err := s.bookingRepo.Create(booking); err != nil {
		return nil, errors.NewDatabaseError("create booking", err)
	}

	// 10. 加载关联的房间信息
	booking.Room = *room

	return booking, nil
}

// GetBookingByID 根据 ID 获取预订详情
func (s *BookingService) GetBookingByID(id int64, userID int64) (*models.Booking, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("预订不存在")
		}
		return nil, errors.NewDatabaseError("find booking", err)
	}

	// 权限检查：只能查看自己的订单（管理员除外）
	// 注意：这里简化处理，实际应该检查用户角色
	if booking.UserID.Int64() != userID {
		return nil, errors.NewForbiddenError("无权访问此预订")
	}

	return booking, nil
}

// GetByGuestInfo 通过客人姓名、手机号和状态查询预订
func (s *BookingService) GetByGuestInfo(guestName, guestPhone, status string) ([]models.Booking, error) {
	// 参数验证
	if guestName == "" && guestPhone == "" {
		return nil, errors.NewBadRequestError("客人姓名和手机号不能同时为空")
	}

	bookings, err := s.bookingRepo.FindByGuestInfo(guestName, guestPhone, status)
	if err != nil {
		return nil, errors.NewDatabaseError("find by guest info", err)
	}
	return bookings, nil
}

// GetMyBookings 获取我的预订列表
func (s *BookingService) GetMyBookings(userID int64, page, pageSize int) ([]models.Booking, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	bookings, total, err := s.bookingRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("find user bookings", err)
	}

	return bookings, total, nil
}

// CancelBooking 取消预订
func (s *BookingService) CancelBooking(id int64, userID int64, reason string) error {
	// 1. 查找预订
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("预订不存在")
		}
		return errors.NewDatabaseError("find booking", err)
	}

	// 2. 权限检查
	if booking.UserID.Int64() != userID {
		return errors.NewForbiddenError("无权取消此预订")
	}

	// 3. 检查是否可以取消
	if !booking.CanCancel() {
		return errors.NewBadRequestError("该预订无法取消")
	}

	// 4. 更新预订状态
	booking.Status = "cancelled"
	booking.CancelReason = reason

	if err := s.bookingRepo.Update(booking); err != nil {
		return errors.NewDatabaseError("cancel booking", err)
	}

	return nil
}

func (s *BookingService) ConfirmBooking(id int64) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("预订不存在")
		}
		return errors.NewDatabaseError("find booking", err)
	}

	if !booking.IsPending() {
		return errors.NewBadRequestError("只能确认待处理的预订")
	}

	if err := s.bookingRepo.UpdateStatus(id, "confirmed"); err != nil {
		return errors.NewDatabaseError("confirm booking", err)
	}

	return nil
}

// CheckIn 办理入住（管理员）
func (s *BookingService) CheckIn(id int64) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("预订不存在")
		}
		return errors.NewDatabaseError("find booking", err)
	}

	if !booking.CanCheckIn() {
		return errors.NewBadRequestError("该预订无法办理入住")
	}

	// 更新预订状态为入住中
	if err := s.bookingRepo.UpdateStatus(id, "checkin"); err != nil {
		return errors.NewDatabaseError("check in", err)
	}

	// 更新房间状态为已占用
	if err := s.roomRepo.UpdateStatus(uint(booking.RoomID), "occupied"); err != nil {
		return errors.NewDatabaseError("update room status", err)
	}

	return nil
}

// CheckOut 办理退房（管理员）
func (s *BookingService) CheckOut(id int64) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("预订不存在")
		}
		return errors.NewDatabaseError("find booking", err)
	}

	if booking.Status != "checkin" {
		return errors.NewBadRequestError("只能为入住中的订单办理退房")
	}

	// 更新预订状态为已退房
	if err := s.bookingRepo.UpdateStatus(id, "checkout"); err != nil {
		return errors.NewDatabaseError("check out", err)
	}

	// 更新房间状态为可用
	if err := s.roomRepo.UpdateStatus(uint(booking.RoomID), "available"); err != nil {
		return errors.NewDatabaseError("update room status", err)
	}

	return nil
}

// ListAllBookings 获取所有预订列表（管理员）
func (s *BookingService) ListAllBookings(page, pageSize int) ([]models.Booking, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	bookings, total, err := s.bookingRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("list bookings", err)
	}

	return bookings, total, nil
}

// GetBookingsByRoomNumberAndStatus 根据房间号和状态获取预订列表（管理员）
func (s *BookingService) GetBookingsByRoomNumberAndStatus(roomNumber string, status string) ([]models.Booking, error) {
	bookings, err := s.bookingRepo.FindByRoomNumberAndStatus(roomNumber, status)
	if err != nil {
		return nil, errors.NewDatabaseError("find bookings by room number and status", err)
	}

	return bookings, nil
}
