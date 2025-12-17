package models

import (
	"time"
)

// Banner 活动横幅模型
// 对应数据库中的 banners 表
type Banner struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`     // 主键（自增ID）
	Title     string     `gorm:"not null;size:100" json:"title"`         // 标题
	Subtitle  *string    `gorm:"size:255" json:"subtitle"`               // 副标题（可为空）
	ImageURL  string     `gorm:"not null;size:500" json:"image_url"`     // 图片URL
	LinkURL   *string    `gorm:"size:500" json:"link_url"`               // 跳转链接（可为空）
	Status    string     `gorm:"default:'active';size:20" json:"status"` // 状态：active, inactive
	Sort      int        `gorm:"default:0" json:"sort"`                  // 排序，数字越小越靠前
	StartTime *time.Time `json:"start_time"`                             // 活动开始时间（可为空）
	EndTime   *time.Time `json:"end_time"`                               // 活动结束时间（可为空）
	CreatedAt time.Time  `json:"created_at"`                             // 创建时间
	UpdatedAt time.Time  `json:"updated_at"`                             // 更新时间
}

// TableName 指定表名
func (Banner) TableName() string {
	return "banners"
}
