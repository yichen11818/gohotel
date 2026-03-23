package handler

import (
	"gohotel/internal/service"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WorkOrderHandler struct {
	service service.WorkOrderService
}

func NewWorkOrderHandler(service service.WorkOrderService) *WorkOrderHandler {
	return &WorkOrderHandler{service: service}
}

// CreateRepairRequest 创建报修
// @Summary 创建维修请求
// @Description 创建新的维修请求
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param request body object true "创建维修请求 {room_id: int, type: string, description: string}"
// @Success 200 {object} utils.Response "创建成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/work-orders/repair [post]
// @Security Bearer
func (h *WorkOrderHandler) CreateRepairRequest(c *gin.Context) {
	var req struct {
		RoomID      int64  `json:"room_id" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.service.CreateRepairRequest(c.Request.Context(), req.RoomID, req.Type, req.Description, userID.(int64)); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "报修申请提交成功", nil)
}

// CompleteRepair 完成维修
// @Summary 完成维修任务
// @Description 标记维修任务为已完成
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param id path int true "维修任务ID"
// @Param request body object false "备注 {remark: string}"
// @Success 200 {object} utils.Response "操作成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/work-orders/repair/{id}/complete [post]
// @Security Bearer
func (h *WorkOrderHandler) CompleteRepair(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	userID, _ := c.Get("user_id")
	if err := h.service.CompleteRepair(c.Request.Context(), uint(id), userID.(int64), req.Remark); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "维修任务已完成", nil)
}

// ListMaintenance 获取维修列表
// @Summary 获取维修任务列表
// @Description 获取所有维修任务，支持按状态筛选
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param status query string false "状态筛选"
// @Success 200 {object} utils.Response{data=service.MaintenanceListResponse} "获取成功"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/work-orders/repairs [get]
// @Security Bearer
func (h *WorkOrderHandler) ListMaintenance(c *gin.Context) {
	status := c.Query("status")
	list, err := h.service.ListMaintenance(c.Request.Context(), status)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, list)
}

// CreateCleaningTask 创建清洁任务
// @Summary 创建清洁任务
// @Description 创建新的清洁任务
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param request body object true "创建清洁任务 {room_id: int, type: string}"
// @Success 200 {object} utils.Response "创建成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/work-orders/cleaning [post]
// @Security Bearer
func (h *WorkOrderHandler) CreateCleaningTask(c *gin.Context) {
	var req struct {
		RoomID int64  `json:"room_id" binding:"required"`
		Type   string `json:"type" binding:"required"` // daily, checkout, deep
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	if err := h.service.CreateCleaningTask(c.Request.Context(), req.RoomID, req.Type); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "清洁任务已创建", nil)
}

// AssignStaff 指派人员
// @Summary 分派员工
// @Description 为工单分派负责员工
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param id path int true "工单ID"
// @Param request body object true "分派请求 {staff_id: int}"
// @Success 200 {object} utils.Response "分派成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/work-orders/cleaning/{id}/assign [post]
// @Security Bearer
func (h *WorkOrderHandler) AssignStaff(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		StaffID int64 `json:"staff_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	if err := h.service.AssignStaff(c.Request.Context(), uint(id), req.StaffID); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "任务已指派", nil)
}

// CompleteCleaning 完成清洁
// @Summary 完成清洁任务
// @Description 标记清洁任务为已完成
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param id path int true "清洁任务ID"
// @Success 200 {object} utils.Response "操作成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/work-orders/cleaning/{id}/complete [post]
// @Security Bearer
func (h *WorkOrderHandler) CompleteCleaning(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.service.CompleteCleaning(c.Request.Context(), uint(id)); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "清洁任务已完成", nil)
}

// ListHousekeeping 获取清洁任务列表
// @Summary 获取清洁任务列表
// @Description 获取所有清洁任务，支持按状态筛选
// @Tags 工单管理
// @Accept json
// @Produce json
// @Param status query string false "状态筛选"
// @Success 200 {object} utils.Response{data=service.HousekeepingListResponse} "获取成功"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/work-orders/cleanings [get]
// @Security Bearer
func (h *WorkOrderHandler) ListHousekeeping(c *gin.Context) {
	status := c.Query("status")
	list, err := h.service.ListHousekeeping(c.Request.Context(), status)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, list)
}
