package service

import (
	"encoding/json"
	"fmt"
	"gohotel/internal/config"
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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
	Username string `json:"username" binding:"required,min=2,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name"`
}

// AddUserRequest 添加用户请求结构
type AddUserRequest struct {
	Username string `json:"username" binding:"required,min=2,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name"`
	Role     string `json:"role"`
}

// UpdateUserRequest 管理员更新用户请求结构
type UpdateUserRequest struct {
	Username   string   `json:"username" binding:"required,min=2,max=20"`
	Email      string   `json:"email" binding:"required,email"`
	Phone      *string  `json:"phone"`
	RealName   *string  `json:"real_name"`
	Avatar     *string  `json:"avatar"`
	Role       *string  `json:"role"`
	Status     *string  `json:"status"`
	Level      *string  `json:"level"`
	Points     *int     `json:"points"`
	Balance    *float64 `json:"balance"`
	TotalSpend *float64 `json:"total_spend"`
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

// UserPayload 通用用户载荷
type UserPayload struct {
	User *models.User `json:"user"`
}

// WeChatLoginRequest 微信登录请求结构
type WeChatLoginRequest struct {
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// WeChatSessionResponse 微信 Code2Session 响应
type WeChatSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
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
	account := strings.TrimSpace(req.Username)
	if account == "" {
		return nil, errors.NewUnauthorizedError("用户名或密码错误")
	}

	var (
		user *models.User
		err  error
	)

	// 允许使用 邮箱/手机号/用户名 登录
	// 说明：这里返回统一的“用户名或密码错误”，避免暴露账号是否存在
	if strings.Contains(account, "@") {
		user, err = s.userRepo.FindByEmail(account)
	} else if regexp.MustCompile(`^\d{6,20}$`).MatchString(account) {
		user, err = s.userRepo.FindByPhone(account)
	} else {
		user, err = s.userRepo.FindByUsername(account)
	}
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
	token, err := utils.GenerateToken(user.ID.Int64(), user.Username, user.Role, user.TokenVersion)
	if err != nil {
		return nil, errors.NewInternalServerError("生成令牌失败")
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

// WeChatLogin 微信一键登录
func (s *UserService) WeChatLogin(req *WeChatLoginRequest) (*LoginResponse, error) {
	// 1. 获取微信 OpenID
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.AppConfig.WeChat.AppID,
		config.AppConfig.WeChat.AppSecret,
		req.Code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.NewInternalServerError("请求微信接口失败")
	}
	defer resp.Body.Close()

	var wxResp WeChatSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&wxResp); err != nil {
		return nil, errors.NewInternalServerError("解析微信响应失败")
	}

	if wxResp.ErrCode != 0 {
		return nil, errors.NewBadRequestError(fmt.Sprintf("微信登录失败: %s", wxResp.ErrMsg))
	}

	// 2. 查找或创建用户
	user, err := s.userRepo.FindByOpenID(wxResp.OpenID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.NewDatabaseError("find user by openid", err)
	}

	if user == nil {
		// 用户不存在，自动注册
		userID := utils.GenID()
		username := "wx_" + utils.GenerateRandomString(8)
		// 确保用户名唯一
		for {
			exists, _ := s.userRepo.ExistsByUsername(username)
			if !exists {
				break
			}
			username = "wx_" + utils.GenerateRandomString(8)
		}

		user = &models.User{
			ID:         utils.JSONInt64(userID),
			Username:   username,
			Email:      username + "@wechat.com", // 生成虚拟邮箱
			Password:   "",                       // 微信登录无密码
			RealName:   req.Nickname,
			Avatar:     req.Avatar,
			Role:       "user",
			Status:     "active",
			OpenID:     &wxResp.OpenID,
			FirstLogin: true,
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, errors.NewDatabaseError("create wechat user", err)
		}
	} else {
		// 用户已存在，更新信息（可选）
		// if req.Nickname != "" && user.RealName == "" {
		// 	user.RealName = req.Nickname
		// 	s.userRepo.Update(user)
		// }
	}

	// 3. 检查账号状态
	if !user.IsActive() {
		return nil, errors.NewForbiddenError("账号已被封禁")
	}

	// 4. 生成 JWT 令牌
	token, err := utils.GenerateToken(user.ID.Int64(), user.Username, user.Role, user.TokenVersion)
	if err != nil {
		return nil, errors.NewInternalServerError("生成令牌失败")
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID int64, phone, realName, avatar *string) (*models.User, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("用户不存在")
		}
		return nil, errors.NewDatabaseError("find user", err)
	}

	// 2. 检查手机号是否已被其他用户使用（如果要更新手机号）
	if phone != nil {
		trimmed := strings.TrimSpace(*phone)
		if trimmed == "" {
			user.Phone = nil
		} else {
			var currentPhone string
			if user.Phone != nil {
				currentPhone = strings.TrimSpace(*user.Phone)
			}
			if trimmed != currentPhone {
				exists, err := s.userRepo.ExistsByPhoneExcludingUser(trimmed, userID)
				if err != nil {
					return nil, errors.NewDatabaseError("check phone", err)
				}
				if exists {
					return nil, errors.NewConflictError("手机号已被使用")
				}
			}
			phoneCopy := trimmed
			user.Phone = &phoneCopy
		}
	}
	if realName != nil {
		user.RealName = strings.TrimSpace(*realName)
	}
	if avatar != nil {
		user.Avatar = strings.TrimSpace(*avatar)
	}

	// 3. 保存更新
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.NewDatabaseError("update user", err)
	}

	updatedUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.NewDatabaseError("reload user", err)
	}

	return updatedUser, nil
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
	user.TokenVersion++
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
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.RealName = strings.TrimSpace(req.RealName)
	req.Role = strings.TrimSpace(req.Role)

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

	if req.Phone != "" {
		exists, err = s.userRepo.ExistsByPhone(req.Phone)
		if err != nil {
			return nil, errors.NewDatabaseError("check phone", err)
		}
		if exists {
			return nil, errors.NewConflictError("手机号已被使用")
		}
	}

	if req.Role == "" {
		req.Role = "user"
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
		Role:       req.Role,
		Status:     "active",
		FirstLogin: true,
	}
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

// UpdateUser 管理员更新用户
func (s *UserService) UpdateUser(userID int64, req *UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("用户不存在")
		}
		return nil, errors.NewDatabaseError("find user", err)
	}

	exists, err := s.userRepo.ExistsByUsernameExcludingUser(req.Username, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("check username", err)
	}
	if exists {
		return nil, errors.NewConflictError("用户名已存在")
	}

	exists, err = s.userRepo.ExistsByEmailExcludingUser(req.Email, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("check email", err)
	}
	if exists {
		return nil, errors.NewConflictError("邮箱已被使用")
	}

	user.Username = strings.TrimSpace(req.Username)
	user.Email = strings.TrimSpace(req.Email)

	if req.RealName != nil {
		user.RealName = strings.TrimSpace(*req.RealName)
	}
	if req.Avatar != nil {
		user.Avatar = strings.TrimSpace(*req.Avatar)
	}

	tokenBump := false
	if req.Role != nil {
		roleVal := strings.TrimSpace(*req.Role)
		if roleVal != "" && roleVal != user.Role {
			user.Role = roleVal
			tokenBump = true
		}
	}
	if req.Status != nil {
		statusVal := strings.TrimSpace(*req.Status)
		if statusVal != "" && statusVal != user.Status {
			user.Status = statusVal
			tokenBump = true
		}
	}
	if req.Level != nil {
		levelVal := strings.TrimSpace(*req.Level)
		if levelVal != "" {
			user.Level = levelVal
		}
	}
	if req.Points != nil {
		user.Points = *req.Points
	}
	if req.Balance != nil {
		user.Balance = *req.Balance
	}
	if req.TotalSpend != nil {
		user.TotalSpend = *req.TotalSpend
	}

	if req.Phone != nil {
		trimmed := strings.TrimSpace(*req.Phone)
		if trimmed == "" {
			user.Phone = nil
		} else {
			var currentPhone string
			if user.Phone != nil {
				currentPhone = strings.TrimSpace(*user.Phone)
			}
			if trimmed != currentPhone {
				exists, err = s.userRepo.ExistsByPhoneExcludingUser(trimmed, userID)
				if err != nil {
					return nil, errors.NewDatabaseError("check phone", err)
				}
				if exists {
					return nil, errors.NewConflictError("手机号已被使用")
				}
			}
			phoneCopy := trimmed
			user.Phone = &phoneCopy
		}
	}

	if tokenBump {
		user.TokenVersion++
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.NewDatabaseError("update user", err)
	}

	updatedUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.NewDatabaseError("reload user", err)
	}

	return updatedUser, nil
}

// UpdateSpendAndPoints 更新消费金额和积分
func (s *UserService) UpdateSpendAndPoints(userID int64, spend float64, points int) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.TotalSpend += spend
	user.Points += points

	// 简单的等级提升逻辑
	if user.TotalSpend >= 10000 {
		user.Level = "platinum"
	} else if user.TotalSpend >= 5000 {
		user.Level = "gold"
	} else if user.TotalSpend >= 2000 {
		user.Level = "silver"
	}

	return s.userRepo.Update(user)
}

// DeleteUsers 批量删除用户
func (s *UserService) DeleteUsers(operatorUserID int64, req *DeleteUsersRequest) error {
	// 检查用户ID列表是否为空
	if len(req.UserIDs) == 0 {
		return errors.NewBadRequestError("用户ID列表不能为空")
	}

	targetAdminCount := int64(0)
	normalizedUserIDs := make([]string, 0, len(req.UserIDs))
	seen := make(map[string]struct{}, len(req.UserIDs))

	for _, rawID := range req.UserIDs {
		userID := strings.TrimSpace(rawID)
		if userID == "" {
			return errors.NewBadRequestError("存在无效的用户ID")
		}
		if _, exists := seen[userID]; exists {
			continue
		}
		seen[userID] = struct{}{}

		parsedID, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			return errors.NewBadRequestError("存在无效的用户ID")
		}
		if parsedID == operatorUserID {
			return errors.NewBadRequestError("不能删除当前登录账号")
		}

		user, err := s.userRepo.FindByID(parsedID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.NewNotFoundError("待删除用户不存在")
			}
			return errors.NewDatabaseError("find user", err)
		}
		if user.Role == "admin" {
			targetAdminCount++
		}

		normalizedUserIDs = append(normalizedUserIDs, userID)
	}

	if len(normalizedUserIDs) == 0 {
		return errors.NewBadRequestError("用户ID列表不能为空")
	}

	if targetAdminCount > 0 {
		adminCount, err := s.userRepo.CountByRole("admin")
		if err != nil {
			return errors.NewDatabaseError("count admins", err)
		}
		if adminCount-targetAdminCount < 1 {
			return errors.NewBadRequestError("至少保留一个管理员账号")
		}
	}

	// 执行批量删除操作
	err := s.userRepo.BatchDelete(normalizedUserIDs)
	if err != nil {
		return errors.NewDatabaseError("batch delete users", err)
	}

	return nil
}
