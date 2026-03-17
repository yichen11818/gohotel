package handler

import (
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type PricingHandler struct {
	repo repository.PricingRepository
}

func NewPricingHandler(repo repository.PricingRepository) *PricingHandler {
	return &PricingHandler{repo: repo}
}

// CreateRule 创建定价规则
// @Summary 创建定价规则
// @Description 创建新的定价规则
// @Tags 定价管理
// @Accept json
// @Produce json
// @Param request body object true "创建定价规则请求 {name: string, type: string, room_type: string, start_date: string, end_date: string, adjustment: number, is_percent: boolean, priority: int}"
// @Success 200 {object} utils.Response "创建成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/pricing/rules [post]
// @Security Bearer
func (h *PricingHandler) CreateRule(c *gin.Context) {
	var req struct {
		Name       string  `json:"name" binding:"required"`
		Type       string  `json:"type" binding:"required"` // holiday, weekend, special
		RoomType   string  `json:"room_type"`
		StartDate  string  `json:"start_date" binding:"required"`
		EndDate    string  `json:"end_date" binding:"required"`
		Adjustment float64 `json:"adjustment" binding:"required"`
		IsPercent  bool    `json:"is_percent"`
		Priority   int     `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	rule := &models.PricingRule{
		Name:       req.Name,
		Type:       req.Type,
		RoomType:   req.RoomType,
		StartDate:  startDate,
		EndDate:    endDate,
		Adjustment: req.Adjustment,
		IsPercent:  req.IsPercent,
		Priority:   req.Priority,
	}

	if err := h.repo.CreateRule(c.Request.Context(), rule); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "定价规则创建成功", rule)
}

// ListRules 获取所有定价规则
// @Summary 获取定价规则列表
// @Description 获取所有定价规则
// @Tags 定价管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "获取成功"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/admin/pricing/rules [get]
// @Security Bearer
func (h *PricingHandler) ListRules(c *gin.Context) {
	rules, err := h.repo.ListRules(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessResponse(c, rules)
}
