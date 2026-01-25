package handler

import (
	"gohotel/internal/service"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HotelSettingsHandler struct {
	hsService *service.HotelSettingsService
}

func NewHotelSettingsHandler(hsService *service.HotelSettingsService) *HotelSettingsHandler {
	return &HotelSettingsHandler{hsService: hsService}
}

func (h *HotelSettingsHandler) GetPublicSettings(c *gin.Context) {
	hotelID, err := strconv.ParseInt(c.Query("hotel_id"), 10, 64)
	if err != nil || hotelID <= 0 {
		utils.ErrorResponse(c, errors.NewBadRequestError("hotel_id 必须为有效数字"))
		return
	}

	settings, err := h.hsService.GetPublicSettingsByHotelID(hotelID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, settings)
}

func (h *HotelSettingsHandler) GetAdminSettings(c *gin.Context) {
	hotelID, err := strconv.ParseInt(c.Query("hotel_id"), 10, 64)
	if err != nil || hotelID <= 0 {
		utils.ErrorResponse(c, errors.NewBadRequestError("hotel_id 必须为有效数字"))
		return
	}

	settings, err := h.hsService.GetSettingsByHotelID(hotelID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, settings)
}

func (h *HotelSettingsHandler) SaveSettings(c *gin.Context) {
	var req service.SaveHotelSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}
	if err := h.hsService.Save(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessWithMessage(c, "保存成功", nil)
}
