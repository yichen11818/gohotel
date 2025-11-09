package service

import (
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"

	"gorm.io/gorm"
)

// UserService 用户业务逻辑层
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// Register 用户注册
func (s *UserService) Register(req *RegisterRequest) (*models.User, error) {
	// 1. 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, errors.NewDatabaseError("check username", err)
	}
	if exists {
		return nil, errors.NewConflictError("用户名已存在")
	}

	// 2. 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, errors.NewDatabaseError("check email", err)
	}
	if exists {
		return nil, errors.NewConflictError("邮箱已被使用")
	}

	// 3. 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalServerError("密码加密失败")
	}

	// 4. 创建用户对象
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
		RealName: req.RealName,
		Role:     "user",
		Status:   "active",
	}

	// 5. 保存到数据库
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.NewDatabaseError("create user", err)
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewUnauthorizedError("用户名或密码错误")
		}
		return nil, errors.NewDatabaseError("find user", err)
	}

	// 2. 检查账号状态
	if !user.IsActive() {
		return nil, errors.NewForbiddenError("账号已被封禁")
	}

	// 3. 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.NewUnauthorizedError("用户名或密码错误")
	}

	// 4. 生成 JWT 令牌
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.NewInternalServerError("生成令牌失败")
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

// GetUserByID 根据 ID 获取用户信息
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("用户不存在")
		}
		return nil, errors.NewDatabaseError("find user", err)
	}
	return user, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID uint, phone, realName, avatar string) (*models.User, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("用户不存在")
		}
		return nil, errors.NewDatabaseError("find user", err)
	}

	// 2. 更新字段
	if phone != "" {
		user.Phone = phone
	}
	if realName != "" {
		user.RealName = realName
	}
	if avatar != "" {
		user.Avatar = avatar
	}

	// 3. 保存更新
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.NewDatabaseError("update user", err)
	}

	return user, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// 1. 查找用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("用户不存在")
		}
		return errors.NewDatabaseError("find user", err)
	}

	// 2. 验证旧密码
	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.NewBadRequestError("旧密码错误")
	}

	// 3. 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.NewInternalServerError("密码加密失败")
	}

	// 4. 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return errors.NewDatabaseError("update password", err)
	}

	return nil
}

// ListUsers 获取用户列表（分页）
func (s *UserService) ListUsers(page, pageSize int) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.userRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("list users", err)
	}

	return users, total, nil
}
