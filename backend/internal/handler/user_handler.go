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
// @Summary 用户注册
// @Description 新用户注册接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 200 {object} models.User
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /api/auth/register [post]
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
// @Summary 用户登录
// @Description 用户登录接口，返回用户信息和 JWT token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/auth/login [post]
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
// @Summary 获取个人信息
// @Description 获取当前登录用户的个人信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} models.User
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取当前用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, errors.NewUnauthorizedError("未登录"))
		return
	}

	user, err := h.userService.GetUserByID(userID.(int64))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, user)
}

// UpdateProfile 更新个人信息
// @Summary 更新个人信息
// @Description 更新当前登录用户的个人信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body object true "更新信息" example({"phone":"13800138000","real_name":"张三","avatar":"https://example.com/avatar.jpg"})
// @Success 200 {object} models.User
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/users/profile [post]
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

	user, err := h.userService.UpdateProfile(userID.(int64), req.Phone, req.RealName, req.Avatar)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "更新成功", user)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body object true "密码信息" example({"old_password":"old123456","new_password":"new123456"})
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/users/password [post]
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

	err := h.userService.ChangePassword(userID.(int64), req.OldPassword, req.NewPassword)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithMessage(c, "密码修改成功", nil)
}

// GetUserByID 根据 ID 获取用户（管理员）
// @Summary 获取用户详情（管理员）
// @Description 管理员根据用户ID获取用户详细信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户 ID"
// @Success 200 {object} models.User
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/admin/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, errors.NewBadRequestError("无效的用户ID"))
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, user)
}

// ListUsers 获取用户列表（管理员）
// @Summary 获取用户列表（管理员）
// @Description 管理员获取所有用户列表，支持分页
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param username query string false "用户名"
// @Param email query string false "邮箱"
// @Param phone query string false "手机号"
// @Param real_name query string false "真实姓名"
// @Param role query string false "角色"
// @Param status query string false "状态"
// @Success 200 {array} models.User
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Router /api/admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	username := c.Query("username")
	email := c.Query("email")
	phone := c.Query("phone")
	realName := c.Query("real_name")
	role := c.Query("role")
	status := c.Query("status")

	users, total, err := h.userService.GetUser(page, pageSize, username, email, phone, realName, role, status)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessWithPage(c, users, page, pageSize, total)
}

// AddUser 添加用户
// @Summary 添加用户
// @Description 管理员添加新的用户账户，默认密码为yumi123456
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.AddUserRequest true "管理员信息"
// @Success 200 {object} models.User
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /api/admin/users/user [post]
func (h *UserHandler) AddUser(c *gin.Context) {
	// 1. 绑定并验证请求参数
	var req service.AddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 2. 调用 Service 层
	user, err := h.userService.AddUser(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 3. 返回成功响应
	utils.SuccessWithMessage(c, "用户添加成功", gin.H{
		"user": user,
	})
}

// DeleteUsers 批量删除用户
// @Summary 批量删除用户
// @Description 管理员批量删除用户账户
// @Tags 管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body service.DeleteUsersRequest true "删除用户信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Router /api/admin/users/batch [post]
func (h *UserHandler) DeleteUsers(c *gin.Context) {
	// 1. 绑定并验证请求参数
	var req service.DeleteUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 2. 调用 Service 层
	err := h.userService.DeleteUsers(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 3. 返回成功响应
	utils.SuccessWithMessage(c, "用户删除成功", nil)
}
