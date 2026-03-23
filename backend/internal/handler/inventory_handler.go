package handler

import (
	"gohotel/internal/service"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// InitInventory 初始化库存
// @Summary 初始化库存
// @Description 初始化房型库存
// @Tags 库存管理
// @Accept json
// @Produce json
// @Param request body object true "初始化库存请求 {room_type: string, total_count: int, price: number, days: int}"
// @Success 200 {object} utils.Response "初始化成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/inventory/init [post]
// @Security Bearer
func (h *InventoryHandler) InitInventory(c *gin.Context) {
	var req struct {
		RoomType   string  `json:"room_type" binding:"required"`
		TotalCount int     `json:"total_count" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
		Days       int     `json:"days" binding:"required"` // 初始化未来多少天的库存
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	if err := h.service.InitInventory(c.Request.Context(), req.RoomType, req.TotalCount, req.Price, req.Days); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "库存初始化成功", nil)
}

// UpdateInventory 调整库存
// @Summary 更新库存
// @Description 调整指定日期范围的库存
// @Tags 库存管理
// @Accept json
// @Produce json
// @Param request body object true "更新库存请求 {room_type: string, start_date: string, end_date: string, delta: int}"
// @Success 200 {object} utils.Response "更新成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/inventory/update [post]
// @Security Bearer
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	var req struct {
		RoomType  string `json:"room_type" binding:"required"`
		StartDate string `json:"start_date" binding:"required"`
		EndDate   string `json:"end_date" binding:"required"`
		Delta     int    `json:"delta" binding:"required"` // 调整量
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	if err := h.service.UpdateInventory(c.Request.Context(), req.RoomType, startDate, endDate, req.Delta); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "库存调整成功", nil)
}

// GetInventoryGrid 获取房态库存格子
// @Summary 获取库存网格
// @Description 获取指定日期范围的库存网格数据
// @Tags 库存管理
// @Accept json
// @Produce json
// @Param start_date query string true "开始日期"
// @Param end_date query string true "结束日期"
// @Success 200 {object} utils.Response{data=service.InventoryGridResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/inventory/grid [get]
// @Security Bearer
func (h *InventoryHandler) GetInventoryGrid(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.ErrorResponse(c, errors.NewBadRequestError("开始日期和结束日期不能为空"))
		return
	}

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	grid, err := h.service.GetInventoryGrid(c.Request.Context(), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, grid)
}
