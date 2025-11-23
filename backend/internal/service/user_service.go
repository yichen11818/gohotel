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

// AddUserRequest 添加用户请求结构
type AddUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name"`
	Role     string `json:"role"`
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

// DeleteUsersRequest 批量删除用户请求结构
type DeleteUsersRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1"`
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

	// 3. 检查手机号是否已存在（如果提供了手机号）
	if req.Phone != "" {
		exists, err = s.userRepo.ExistsByPhone(req.Phone)
		if err != nil {
			return nil, errors.NewDatabaseError("check phone", err)
		}
		if exists {
			return nil, errors.NewConflictError("手机号已被使用")
		}
	}

	// 4. 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalServerError("密码加密失败")
	}

	// 5. 生成雪花 ID
	userID := utils.GenID()

	// 6. 创建用户对象
	var phonePtr *string
	if req.Phone != "" {
		phonePtr = &req.Phone
	}
	user := &models.User{
		ID:       utils.JSONInt64(userID),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    phonePtr,
		RealName: req.RealName,
		Role:     "user",
		Status:   "active",
	}

	// 7. 保存到数据库
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
	token, err := utils.GenerateToken(user.ID.Int64(), user.Username)
	if err != nil {
		return nil, errors.NewInternalServerError("生成令牌失败")
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID int64, phone, realName, avatar string) (*models.User, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("用户不存在")
		}
		return nil, errors.NewDatabaseError("find user", err)
	}

	// 2. 检查手机号是否已被其他用户使用（如果要更新手机号）
	if phone != "" {
		var currentPhone string
		if user.Phone != nil {
			currentPhone = *user.Phone
		}
		if phone != currentPhone {
			exists, err := s.userRepo.ExistsByPhoneExcludingUser(phone, userID)
			if err != nil {
				return nil, errors.NewDatabaseError("check phone", err)
			}
			if exists {
				return nil, errors.NewConflictError("手机号已被使用")
			}
			phonePtr := &phone
			user.Phone = phonePtr
		}
	} else {
		// 如果传入空字符串，清空手机号
		user.Phone = nil
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
func (s *UserService) ChangePassword(userID int64, oldPassword, newPassword string) error {
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

// GetUserByID 根据 ID 获取用户信息
func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("用户不存在")
		}
		return nil, errors.NewDatabaseError("find user", err)
	}
	return user, nil
}

// GetUser 根据条件查询用户
func (s *UserService) GetUser(page, pageSize int, username, email, phone, realName, role, status string) ([]models.User, int64, error) {

	users, total, err := s.userRepo.FindAll(page, pageSize, username, email, phone, realName, role, status)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("list users filter", err)
	}
	return users, total, nil
}

// AddUser 添加用户
func (s *UserService) AddUser(req *AddUserRequest) (*models.User, error) {
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

	// 3. 加密默认密码
	hashedPassword, err := utils.HashPassword("yumi123456")
	if err != nil {
		return nil, errors.NewInternalServerError("密码加密失败")
	}

	// 4. 生成雪花 ID
	userID := utils.GenID()

	// 5. 创建用户对象
	user := &models.User{
		ID:         utils.JSONInt64(userID),
		Username:   req.Username,
		Email:      req.Email,
		Password:   hashedPassword,
		RealName:   req.RealName,
		Role:       "user", // 默认角色为 user
		Status:     "active",
		FirstLogin: true,
	}
	user.Role = req.Role // 处理 Phone 字段，将 string 转换为 *string
	if req.Phone != "" {
		phone := req.Phone
		user.Phone = &phone
	}
	// 6. 保存到数据库
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.NewDatabaseError("create admin user", err)
	}

	return user, nil
}

// DeleteUsers 批量删除用户
func (s *UserService) DeleteUsers(req *DeleteUsersRequest) error {
	// 检查用户ID列表是否为空
	if len(req.UserIDs) == 0 {
		return errors.NewBadRequestError("用户ID列表不能为空")
	}

	// 执行批量删除操作
	err := s.userRepo.BatchDelete(req.UserIDs)
	if err != nil {
		return errors.NewDatabaseError("batch delete users", err)
	}

	return nil
}
