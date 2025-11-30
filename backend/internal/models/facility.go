package models

import (
	"time"
)

// Facility 设施模型
// 对应数据库中的 facilities 表
type Facility struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Type      string    `gorm:"not null;size:50;index" json:"type"` // 设施类型：elevator, corridor, storage 等
	Floor     int       `gorm:"not null;index" json:"floor"`        // 楼层
	Left      int       `gorm:"not null" json:"left"`               // X 坐标（左边距）
	Top       int       `gorm:"not null" json:"top"`                // Y 坐标（上边距）
	Width     int       `gorm:"not null" json:"width"`              // 宽度
	Height    int       `gorm:"not null" json:"height"`             // 高度
	Rotation  int       `gorm:"default:0" json:"rotation"`          // 旋转角度
	Label     string    `gorm:"size:100" json:"label"`              // 设施标签（可选）
	CreatedAt time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                         // 更新时间
}

// TableName 指定表名
func (Facility) TableName() string {
	return "facilities"
}
