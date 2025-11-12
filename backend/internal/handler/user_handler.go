package handler

import (
	"gohotel/internal/service"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户控制器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户控制器实例
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register 用户注册
// POST /api/auth/register
func (h *UserHandler) Register(c *gin.Context) {

	var req service.RegisterRequest
	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 2. 调用 Service 层
	user, err := h.userService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 3. 返回成功响应
	utils.SuccessWithMessage(c, "注册成功", gin.H{
		"user": user,
	})
}

// Login 用户登录
// POST /api/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest

	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 2. 调用 Service 层
	resp, err := h.userService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 3. 返回成功响应
	utils.SuccessWithMessage(c, "登录成功", resp)
}

// GetProfile 获取个人信息
// GET /api/users/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取当前用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.NewUnauthorizedError("未登录"))
		return
	}

	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, user)
}

// UpdateProfile 更新个人信息
// POST /api/users/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// 获取当前用户 ID
	userID, _ := c.Get("user_id")

	var req struct {
		Phone    string `json:"phone"`
		RealName string `json:"real_name"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	user, err := h.userService.UpdateProfile(userID.(uint), req.Phone, req.RealName, req.Avatar)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "更新成功", user)
}

// ChangePassword 修改密码
// POST /api/users/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	err := h.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "密码修改成功", nil)
}

// GetUserByID 根据 ID 获取用户（管理员）
// GET /api/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的用户ID"))
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, user)
}

// ListUsers 获取用户列表（管理员）
// GET /api/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.userService.ListUsers(page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithPage(c, users, page, pageSize, total)
}
