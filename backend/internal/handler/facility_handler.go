package handler

import (
	"gohotel/internal/service"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FacilityHandler 设施控制器
type FacilityHandler struct {
	facilityService *service.FacilityService
}

// NewFacilityHandler 创建设施控制器实例
func NewFacilityHandler(facilityService *service.FacilityService) *FacilityHandler {
	return &FacilityHandler{facilityService: facilityService}
}

// CreateFacility 创建设施
// @Summary 创建设施
// @Description 创建设施
// @Tags 设施
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CreateFacilityRequest true "设施信息"
// @Success 200 {object} models.Facility
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /api/facilities [post]
func (h *FacilityHandler) CreateFacility(c *gin.Context) {
	var req service.CreateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}
	facility, err := h.facilityService.CreateFacility(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessWithMessage(c, "设施创建成功", facility)
}

// UpdateFacility 更新设施
// @Summary 更新设施
// @Description 更新设施
// @Tags 设施
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设施 ID"
// @Param request body service.UpdateFacilityRequest true "设施信息"
// @Success 200 {object} models.Facility
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/facilities/{id} [put]
func (h *FacilityHandler) UpdateFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的设施ID"))
		return
	}

	var req service.UpdateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	facility, err := h.facilityService.UpdateFacility(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "设施更新成功", facility)
}

// DeleteFacility 删除设施
// @Summary 删除设施
// @Description 删除设施
// @Tags 设施
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "设施 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/facilities/{id} [delete]
func (h *FacilityHandler) DeleteFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的设施ID"))
		return
	}
	err = h.facilityService.DeleteFacility(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessWithMessage(c, "设施删除成功", nil)
}

// FindFacilityByID 根据 ID 查找设施
// @Summary 根据 ID 查找设施
// @Description 根据 ID 查找设施
// @Tags 设施
// @Accept json
// @Produce json
// @Param id path int true "设施 ID"
// @Success 200 {object} models.Facility
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/facilities/{id} [get]
func (h *FacilityHandler) FindFacilityByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的设施ID"))
		return
	}
	facility, err := h.facilityService.FindFacilityByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, facility)
}

// FindAllFacilities 查询所有设施
// @Summary 查询所有设施
// @Description 查询所有设施（分页）
// @Tags 设施
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {array} models.Facility
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/facilities [get]
func (h *FacilityHandler) FindAllFacilities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	facilities, total, err := h.facilityService.FindAllFacilities(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessWithPage(c, facilities, page, pageSize, total)
}

// FindFacilitiesByFloor 按楼层查询设施
// @Summary 按楼层查询设施
// @Description 获取指定楼层的所有设施
// @Tags 设施
// @Accept json
// @Produce json
// @Param floor path int true "楼层"
// @Success 200 {array} models.Facility
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/facilities/floor/{floor} [get]
func (h *FacilityHandler) FindFacilitiesByFloor(c *gin.Context) {
	floor, err := strconv.Atoi(c.Param("floor"))
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的楼层"))
		return
	}
	facilities, err := h.facilityService.FindFacilitiesByFloor(floor)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, facilities)
}

// BatchUpdateFacilities 批量更新设施位置
// @Summary 批量更新设施位置
// @Description 批量更新设施的位置和尺寸信息
// @Tags 设施
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.BatchUpdateFacilitiesRequest true "设施位置信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Router /api/facilities/batch [put]
func (h *FacilityHandler) BatchUpdateFacilities(c *gin.Context) {
	var req service.BatchUpdateFacilitiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}
	if err := h.facilityService.BatchUpdateFacilities(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessWithMessage(c, "设施位置已批量更新", nil)
}
