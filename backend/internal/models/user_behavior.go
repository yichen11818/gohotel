package models

import (
	"gohotel/pkg/utils"
	"time"
)

// UserBehavior 用户行为模型
// 用于记录浏览、点击、收藏意图等隐式反馈，为推荐算法提供数据基础
type UserBehavior struct {
	ID           utils.JSONInt64 `gorm:"primaryKey;autoIncrement:false" json:"id"`
	UserID       utils.JSONInt64 `gorm:"not null;index" json:"user_id"`
	RoomID       int64           `gorm:"not null;index" json:"room_id"`
	BehaviorType string          `gorm:"not null;size:50;index" json:"behavior_type"`
	Source       string          `gorm:"size:50" json:"source"`
	Weight       float64         `gorm:"not null;type:decimal(8,2)" json:"weight"`
	MetadataJSON string          `gorm:"type:text" json:"metadata_json,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`

	Room Room `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}

// TableName 指定表名
func (UserBehavior) TableName() string {
	return "user_behaviors"
}
