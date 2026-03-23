package handler

import (
	"strconv"

	"gohotel/internal/service"
	"gohotel/pkg/utils"

	"github.com/gin-gonic/gin"
)

// LogHandler 日志控制器
type LogHandler struct {
	logService *service.LogService
}

// NewLogHandler 创建日志控制器实例
func NewLogHandler(logService *service.LogService) *LogHandler {
	return &LogHandler{logService: logService}
}

// Report 日志上报接口
// @Summary 日志上报
// @Description 客户端批量上报日志
// @Tags 日志
// @Accept json
// @Produce json
// @Param request body service.LogReportRequest true "日志数据"
// @Success 200 {object} utils.Response
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/logs/report [post]
func (h *LogHandler) Report(c *gin.Context) {
	var req service.LogReportRequest

	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 2. 调用 Service 层批量存储日志
	if err := h.logService.ReportLogs(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 3. 返回成功响应
	utils.SuccessWithMessage(c, "日志上报成功", gin.H{
		"count": len(req.Logs),
	})
}

// GetLogs 获取日志列表
// @Summary 获取日志列表
// @Description 分页获取日志列表（管理员）
// @Tags 日志
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} utils.PageResponse{data=service.LogListResponse}
// @Failure 400 {object} errors.ErrorResponse
// @Security Bearer
// @Router /api/admin/logs [get]
func (h *LogHandler) GetLogs(c *gin.Context) {
	page := getQueryInt(c, "page", 1)
	pageSize := getQueryInt(c, "page_size", 20)

	logs, total, err := h.logService.GetLogs(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithPage(c, logs, page, pageSize, total)
}

// getQueryInt 获取查询参数的整数值
func getQueryInt(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
