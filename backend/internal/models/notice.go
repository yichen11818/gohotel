package models

import (
	"gohotel/pkg/utils"
	"time"
)

// notice 公告模型
// 对应数据库中的 notices表
type Notice struct {
	ID        utils.JSONInt64 `gorm:"primaryKey;json:"id"`            // 自增主键
	Title     string          `gorm:"not null;size:100" json:"title"` // 标题
	LinkURL   *string         `gorm:"size:500" json:"link_url"`       // 跳转链接（可为空）
	Sort      int             `gorm:"default:0" json:"sort"`          // 排序，数字越小越靠前
	Status    string          `gorm:"default:'active';size:20" json:"status"`   // 状态：active, inactive
	StartTime *time.Time      `json:"start_time"`                     // 活动开始时间（可为空）
	EndTime   *time.Time      `json:"end_time"`                       // 活动结束时间（可为空）
	CreatedAt time.Time       `json:"created_at"`                     // 创建时间
	UpdatedAt time.Time       `json:"updated_at"`                     // 更新时间
}
// TableName 指定表名
func (Notice) TableName() string {
	return "notices"
}
