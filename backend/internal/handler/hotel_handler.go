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
