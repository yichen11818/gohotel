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

// GetPublicSettings 获取酒店公开设置
// @Summary 获取酒店公开设置
// @Description 获取指定酒店的公开配置信息
// @Tags 酒店设置
// @Accept json
// @Produce json
// @Param hotel_id query int true "酒店 ID"
// @Success 200 {object} models.HotelSettings
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/settings/public [get]
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

// GetAdminSettings 获取酒店设置（管理员）
// @Summary 获取酒店设置
// @Description 获取指定酒店的配置信息
// @Tags 酒店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param hotel_id query int true "酒店 ID"
// @Success 200 {object} models.HotelSettings
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/admin/settings [get]
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

// SaveSettings 保存酒店设置（管理员）
// @Summary 保存酒店设置
// @Description 保存或更新酒店配置
// @Tags 酒店管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.SaveHotelSettingsRequest true "设置信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/admin/settings/save [post]
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
