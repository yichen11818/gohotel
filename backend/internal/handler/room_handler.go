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
// POST /api/rooms
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
// GET /api/rooms/:id
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
// POST /api/rooms/:id
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
// POST /api/rooms/:id/delete
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
// GET /api/rooms
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
// GET /api/rooms/available
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
// GET /api/rooms/floor/:floor
func (h *RoomHandler) GetRoomByFloor(c *gin.Context) {
	floor,err := strconv.Atoi(c.Param("floor"))//strconv.Atoi 将字符串转换为整数
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
// GET /api/rooms/search/type
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
