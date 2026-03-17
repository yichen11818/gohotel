package handler

import (
	"gohotel/internal/service"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HotelHandler struct {
	hotelService *service.HotelService
}

func NewHotelHandler(hotelService *service.HotelService) *HotelHandler {
	return &HotelHandler{hotelService: hotelService}
}

// CreateHotel 创建酒店（管理员）
// @Summary 创建酒店
// @Description 创建新的酒店
// @Tags 酒店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CreateHotelRequest true "酒店信息"
// @Success 200 {object} models.Hotel
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/admin/hotels [post]
func (h *HotelHandler) CreateHotel(c *gin.Context) {
	var req service.CreateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}
	hotel, err := h.hotelService.CreateHotel(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, hotel)
}

// GetHotelByID 获取酒店详情（管理员）
// @Summary 获取酒店详情
// @Description 根据ID获取酒店详情
// @Tags 酒店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "酒店 ID"
// @Success 200 {object} models.Hotel
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/admin/hotels/{id} [get]
func (h *HotelHandler) GetHotelByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的酒店ID"))
		return
	}
	hotel, err := h.hotelService.GetHotelByID(id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, hotel)
}

// ListHotels 获取酒店列表（管理员）
// @Summary 获取酒店列表
// @Description 获取所有酒店列表，支持分页
// @Tags 酒店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "{\"hotels\": [...], \"total\": ...}"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/admin/hotels [get]
func (h *HotelHandler) ListHotels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSizeStr := c.DefaultQuery("pageSize", "")
	if pageSizeStr == "" {
		pageSizeStr = c.DefaultQuery("page_size", "10")
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	hotels, total, err := h.hotelService.ListHotels(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hotels":    hotels,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateHotel 更新酒店信息（管理员）
// @Summary 更新酒店信息
// @Description 更新指定酒店的信息
// @Tags 酒店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "酒店 ID"
// @Param request body service.UpdateHotelRequest true "更新信息"
// @Success 200 {object} models.Hotel
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/admin/hotels/{id} [post]
func (h *HotelHandler) UpdateHotel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的酒店ID"))
		return
	}

	var req service.UpdateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	hotel, err := h.hotelService.UpdateHotel(id, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, hotel)
}

// DeleteHotel 删除酒店（管理员）
// @Summary 删除酒店
// @Description 删除指定酒店
// @Tags 酒店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "酒店 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/admin/hotels/{id}/delete [post]
func (h *HotelHandler) DeleteHotel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的酒店ID"))
		return
	}
	if err := h.hotelService.DeleteHotel(id); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}
