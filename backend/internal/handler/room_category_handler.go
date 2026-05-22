package handler

import (
	"strconv"

	"gohotel/internal/service"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RoomCategoryHandler 房型分类控制器
type RoomCategoryHandler struct {
	roomCategoryService *service.RoomCategoryService
}

func NewRoomCategoryHandler(roomCategoryService *service.RoomCategoryService) *RoomCategoryHandler {
	return &RoomCategoryHandler{roomCategoryService: roomCategoryService}
}

// CreateRoomCategory 创建房型分类
// @Summary 创建房型分类（管理员）
// @Description 管理员创建房型分类，统一维护房型描述、预览图和设施
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CreateRoomCategoryRequest true "房型分类信息"
// @Success 200 {object} models.RoomCategory
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /api/admin/room-categories [post]
func (h *RoomCategoryHandler) CreateRoomCategory(c *gin.Context) {
	var req service.CreateRoomCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	category, err := h.roomCategoryService.CreateRoomCategory(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "房型分类创建成功", category)
}

// UpdateRoomCategory 更新房型分类
// @Summary 更新房型分类（管理员）
// @Description 管理员更新房型分类，若修改名称会同步更新房间、库存和定价规则中的房型字段
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "房型分类 ID"
// @Param request body service.UpdateRoomCategoryRequest true "房型分类信息"
// @Success 200 {object} models.RoomCategory
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /api/admin/room-categories/{id} [post]
func (h *RoomCategoryHandler) UpdateRoomCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的房型分类 ID"))
		return
	}

	var req service.UpdateRoomCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	category, err := h.roomCategoryService.UpdateRoomCategory(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "房型分类更新成功", category)
}

// DeleteRoomCategory 删除房型分类
// @Summary 删除房型分类（管理员）
// @Description 管理员删除房型分类；如果仍有关联房间、库存或定价规则，则会阻止删除
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "房型分类 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /api/admin/room-categories/{id}/delete [post]
func (h *RoomCategoryHandler) DeleteRoomCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的房型分类 ID"))
		return
	}

	if err := h.roomCategoryService.DeleteRoomCategory(uint(id)); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "房型分类删除成功", nil)
}

// GetRoomCategoryByID 获取房型分类详情
// @Summary 获取房型分类详情（管理员）
// @Description 管理员根据 ID 获取房型分类详情
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "房型分类 ID"
// @Success 200 {object} models.RoomCategory
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/admin/room-categories/{id} [get]
func (h *RoomCategoryHandler) GetRoomCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的房型分类 ID"))
		return
	}

	category, err := h.roomCategoryService.GetRoomCategoryByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, category)
}

// ListRoomCategories 查询房型分类列表
// @Summary 查询房型分类列表（管理员）
// @Description 管理员分页查询房型分类，返回房型名称、描述、预览图、设施和关联房间数
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(50)
// @Success 200 {array} models.RoomCategory
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Router /api/admin/room-categories [get]
func (h *RoomCategoryHandler) ListRoomCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	categories, total, err := h.roomCategoryService.ListRoomCategories(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithPage(c, categories, page, pageSize, total)
}
