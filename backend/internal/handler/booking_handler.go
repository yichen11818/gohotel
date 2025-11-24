package handler

import (
	"gohotel/internal/service"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"net/http"
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
// @Summary 创建预订
// @Description 创建新的房间预订，需要登录
// @Tags 预订
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CreateBookingRequest true "预订信息"
// @Success 200 {object} models.Booking
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	// 获取当前登录用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.NewUnauthorizedError("未登录"))
		return
	}

	var req service.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	booking, err := h.bookingService.CreateBooking(userID.(int64), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "预订创建成功", booking)
}

// GetBookingByID 获取预订详情
// @Summary 获取预订详情
// @Description 根据预订ID获取预订详细信息，只能查看自己的预订
// @Tags 预订
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "预订 ID"
// @Success 200 {object} models.Booking
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/bookings/{id} [get]
func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的预订ID"))
		return
	}

	booking, err := h.bookingService.GetBookingByID(id, userID.(int64))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, booking)
}

// GetMyBookings 获取我的预订列表
// @Summary 获取我的预订列表
// @Description 获取当前登录用户的所有预订列表，支持分页
// @Tags 预订
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {array} models.Booking
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/bookings/my [get]
func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	bookings, total, err := h.bookingService.GetMyBookings(userID.(int64), page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithPage(c, bookings, page, pageSize, total)
}

// CancelBooking 取消预订
// @Summary 取消预订
// @Description 取消指定的预订，只能取消自己的预订
// @Tags 预订
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "预订 ID"
// @Param request body object true "取消原因" example({"reason":"行程变更"})
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Router /api/bookings/{id}/cancel [post]
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的预订ID"))
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	err = h.bookingService.CancelBooking(id, userID.(int64), req.Reason)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "预订已取消", nil)
}

// ConfirmBooking 确认预订（管理员）
// @Summary 确认预订（管理员）
// @Description 管理员确认待处理的预订
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "预订 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/admin/bookings/{id}/confirm [post]
func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	// 直接获取string类型的ID参数，然后转换为int64
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的预订ID，请确保传入有效的数字字符串"))
		return
	}

	err = h.bookingService.ConfirmBooking(id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "预订已确认", nil)
}

// CheckIn 办理入住（管理员）
// @Summary 办理入住（管理员）
// @Description 管理员为已确认的预订办理入住
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "预订 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/bookings/{id}/checkin [post]
func (h *BookingHandler) CheckIn(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的预订ID"))
		return
	}

	err = h.bookingService.CheckIn(id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "入住办理成功", nil)
}

// CheckOut 办理退房（管理员）
// @Summary 办理退房（管理员）
// @Description 管理员为入住中的预订办理退房
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "预订 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/bookings/{id}/checkout [post]
func (h *BookingHandler) CheckOut(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的预订ID"))
		return
	}

	err = h.bookingService.CheckOut(id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "退房办理成功", nil)
}

// ListAllBookings 获取所有预订（管理员）
// @Summary 获取所有预订（管理员）
// @Description 管理员获取所有预订列表，支持分页
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {array} models.Booking
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Router /api/admin/bookings [get]
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

// SearchBookingsByGuestInfo 通过客人信息搜索预订
// @Summary 通过客人信息搜索预订
// @Description 根据客人姓名和手机号搜索预订记录
// @Tags 管理员
// @Accept json
// @Produce json
// @Param guest_name query string false "客人姓名"
// @Param guest_phone query string false "客人手机号"
// @Success 200 {object} map[string]interface{} "{\"data\": [...], \"count\": number}"
// @Failure 400 {object} map[string]string "{\"error\": string}"
// @Router /api/admin/bookings/search [get]
func (h *BookingHandler) SearchBookingsByGuestInfo(c *gin.Context) {
	guestName := c.Query("guest_name")
	guestPhone := c.Query("guest_phone")

	bookings, err := h.bookingService.GetByGuestInfo(guestName, guestPhone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  bookings,
		"count": len(bookings),
	})
}
