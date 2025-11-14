package handler

import (
	"gohotel/internal/service"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoomHandler 房间控制器
type RoomHandler struct {
	roomService *service.RoomService
}

// NewRoomHandler 创建房间控制器实例
func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

// CreateRoom 创建房间（管理员）
// @Summary 创建房间（管理员）
// @Description 管理员创建新房间
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.CreateRoomRequest true "房间信息"
// @Success 200 {object} models.Room
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /api/rooms [post]
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req service.CreateRoomRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	room, err := h.roomService.CreateRoom(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "房间创建成功", room)
}

// GetRoomByID 根据 ID 获取房间
// @Summary 获取房间详情
// @Description 根据房间ID获取房间详细信息
// @Tags 房间
// @Accept json
// @Produce json
// @Param id path int true "房间 ID"
// @Success 200 {object} models.Room
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/rooms/{id} [get]
func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的房间ID"))
		return
	}

	room, err := h.roomService.GetRoomByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, room)
}

// UpdateRoom 更新房间（管理员）
// @Summary 更新房间（管理员）
// @Description 管理员更新房间信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "房间 ID"
// @Param request body service.UpdateRoomRequest true "房间信息"
// @Success 200 {object} models.Room
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/rooms/{id} [post]
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的房间ID"))
		return
	}

	var req service.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError(err.Error()))
		return
	}

	room, err := h.roomService.UpdateRoom(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "房间更新成功", room)
}

// DeleteRoom 删除房间（管理员）
// @Summary 删除房间（管理员）
// @Description 管理员删除房间
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "房间 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/rooms/{id}/delete [post]
func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的房间ID"))
		return
	}

	err = h.roomService.DeleteRoom(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "房间删除成功", nil)
}

// ListRooms 获取房间列表
// @Summary 获取房间列表
// @Description 获取所有房间列表，支持分页
// @Tags 房间
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {array} models.Room
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/rooms [get]
func (h *RoomHandler) ListRooms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	rooms, total, err := h.roomService.ListRooms(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithPage(c, rooms, page, pageSize, total)
}

// ListAvailableRooms 获取可用房间列表
// @Summary 获取可用房间列表
// @Description 获取所有可用状态的房间列表，支持分页
// @Tags 房间
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {array} models.Room
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/rooms/available [get]
func (h *RoomHandler) ListAvailableRooms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	rooms, total, err := h.roomService.ListAvailableRooms(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithPage(c, rooms, page, pageSize, total)
}

// GetRoomByFloor 根据楼层获取房间
// @Summary 根据楼层获取房间
// @Description 根据楼层号获取该楼层的所有房间，支持分页
// @Tags 房间
// @Accept json
// @Produce json
// @Param floor path int true "楼层号"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {array} models.Room
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/rooms/floor/{floor} [get]
func (h *RoomHandler) GetRoomByFloor(c *gin.Context) {
	floor, err := strconv.Atoi(c.Param("floor")) //strconv.Atoi 将字符串转换为整数
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的楼层"))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	rooms, total, err := h.roomService.ListRoomsByFloor(floor, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}
	utils.SuccessWithPage(c, rooms, page, pageSize, total)
}

// SearchRoomsByType 按房型搜索
// @Summary 按房型搜索房间
// @Description 根据房型搜索房间，支持分页
// @Tags 房间
// @Accept json
// @Produce json
// @Param type query string true "房型"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {array} models.Room
// @Failure 400 {object} errors.ErrorResponse
// @Router /api/rooms/search/type [get]
func (h *RoomHandler) SearchRoomsByType(c *gin.Context) {
	roomType := c.Query("type")
	if roomType == "" {
		utils.ErrorResponse(c, errors.NewBadRequestError("请提供房型参数"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	rooms, total, err := h.roomService.SearchRoomsByType(roomType, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithPage(c, rooms, page, pageSize, total)
}
