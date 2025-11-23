package models

import (
	"gohotel/pkg/utils"
	"time"
)

// User 用户模型
// 对应数据库中的 users 表
type User struct {
	ID         utils.JSONInt64 `gorm:"primaryKey;autoIncrement:false" json:"id"` // 主键（使用雪花算法生成，JSON序列化为字符串）
	Username   string    `gorm:"unique;not null;size:50" json:"username"`  // 用户名（唯一）
	Email      string    `gorm:"unique;not null;size:100" json:"email"`    // 邮箱（唯一）
	Password   string    `gorm:"not null;size:255" json:"-"`               // 密码（不返回给前端）
	Phone      *string   `gorm:"unique;size:20" json:"phone"`              // 手机号（唯一，可为空）
	RealName   string    `gorm:"size:50" json:"real_name"`                 // 真实姓名
	Avatar     string    `gorm:"size:255" json:"avatar"`                   // 头像 URL
	Role       string    `gorm:"default:'user';size:20" json:"role"`       // 角色：user, admin
	Status     string    `gorm:"default:'active';size:20" json:"status"`
	FirstLogin bool      `gorm:"default:false" json:"first_login"` // 是否首次登录
	CreatedAt  time.Time `json:"created_at"`                       // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                       // 更新时间
}

// TableName 指定表名
// 默认情况下，GORM 会使用复数形式（users）
// 这个方法可以自定义表名
func (User) TableName() string {
	return "users"
}

// IsAdmin 判断是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsActive 判断账号是否激活
func (u *User) IsActive() bool {
	return u.Status == "active"
}
