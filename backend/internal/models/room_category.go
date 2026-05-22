package models

import "time"

// RoomCategory 房型分类模型
// 对应数据库中的 room_categories 表
type RoomCategory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`                      // 主键
	Name        string    `gorm:"unique;not null;size:50;index" json:"name"` // 房型名称（唯一）
	Description string    `gorm:"type:text" json:"description"`              // 房型描述
	Facilities  string    `gorm:"type:text" json:"facilities"`               // 设施（JSON 数组字符串）
	Images      string    `gorm:"type:text" json:"images"`                   // 预览图（JSON 数组字符串）
	UsageCount  int64     `gorm:"-" json:"usage_count"`                      // 关联房间数量（运行时统计）
	CreatedAt   time.Time `json:"created_at"`                                // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                                // 更新时间
}

func (RoomCategory) TableName() string {
	return "room_categories"
}
