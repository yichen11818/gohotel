package handler

import (
	"gohotel/internal/service"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BookingHandler 预订控制器
type BookingHandler struct {
	bookingService *service.BookingService
}

// NewBookingHandler 创建预订控制器实例
func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

// CreateBooking 创建预订
// POST /api/bookings
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	// 获取当前登录用户 ID
	userID, _ := c.Get("user_id")

	var req service.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	booking, err := h.bookingService.CreateBooking(userID.(uint), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "预订创建成功", booking)
}

// GetBookingByID 获取预订详情
// GET /api/bookings/:id
func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的预订ID"))
		return
	}

	booking, err := h.bookingService.GetBookingByID(uint(id), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, booking)
}

// GetMyBookings 获取我的预订列表
// GET /api/bookings/my
func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	bookings, total, err := h.bookingService.GetMyBookings(userID.(uint), page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithPage(c, bookings, page, pageSize, total)
}

// CancelBooking 取消预订
// POST /api/bookings/:id/cancel
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的预订ID"))
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	err = h.bookingService.CancelBooking(uint(id), userID.(uint), req.Reason)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "预订已取消", nil)
}

// ConfirmBooking 确认预订（管理员）
// POST /api/bookings/:id/confirm
func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的预订ID"))
		return
	}

	err = h.bookingService.ConfirmBooking(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "预订已确认", nil)
}

// CheckIn 办理入住（管理员）
// POST /api/bookings/:id/checkin
func (h *BookingHandler) CheckIn(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的预订ID"))
		return
	}

	err = h.bookingService.CheckIn(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "入住办理成功", nil)
}

// CheckOut 办理退房（管理员）
// POST /api/bookings/:id/checkout
func (h *BookingHandler) CheckOut(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的预订ID"))
		return
	}

	err = h.bookingService.CheckOut(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "退房办理成功", nil)
}

// ListAllBookings 获取所有预订（管理员）
// GET /api/admin/bookings
func (h *BookingHandler) ListAllBookings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	bookings, total, err := h.bookingService.ListAllBookings(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithPage(c, bookings, page, pageSize, total)
}
