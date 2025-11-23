package repository

import (
	"gohotel/internal/models"
	"strconv"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// ExistsByUsername 检查用户名是否已存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否已存在
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByPhone 检查手机号是否已存在
func (r *UserRepository) ExistsByPhone(phone string) (bool, error) {
	if phone == "" {
		return false, nil
	}
	var count int64
	err := r.db.Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error
	return count > 0, err
}

// ExistsByPhoneExcludingUser 检查手机号是否已被其他用户使用
func (r *UserRepository) ExistsByPhoneExcludingUser(phone string, excludeUserID int64) (bool, error) {
	if phone == "" {
		return false, nil
	}
	var count int64
	err := r.db.Model(&models.User{}).Where("phone = ? AND id != ?", phone, excludeUserID).Count(&count).Error
	return count > 0, err
}

// Update 更新用户信息
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id int64) error {
	return r.db.Delete(&models.User{}, id).Error
}

// BatchDelete 批量删除用户
func (r *UserRepository) BatchDelete(userIDs []string) error {
	// 将字符串ID转换为int64
	var ids []int64
	for _, idStr := range userIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	return r.db.Delete(&models.User{}, ids).Error
}

// FindByID 根据 ID 查找用户
func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll 根据条件过滤查询用户（分页）
func (r *UserRepository) FindAll(page, pageSize int, username, email, phone, realName, role, status string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 构建查询条件
	query := r.db.Model(&models.User{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if realName != "" {
		query = query.Where("real_name LIKE ?", "%"+realName+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 先查询符合条件的总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 计算偏移量并分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}
