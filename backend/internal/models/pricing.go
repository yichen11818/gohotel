package models

import (
	"time"
)

// PricingRule 动态定价规则
type PricingRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;size:100" json:"name"`
	Type        string    `gorm:"not null;size:20" json:"type"`              // 类型：holiday, weekend, special
	RoomType    string    `gorm:"size:50;index" json:"room_type"`            // 适用房型 (空表示全部)
	StartDate   time.Time `gorm:"not null;index" json:"start_date"`
	EndDate     time.Time `gorm:"not null;index" json:"end_date"`
	Adjustment  float64   `gorm:"not null;type:decimal(10,2)" json:"adjustment"` // 调整金额 (+ 为涨价, - 为降价)
	IsPercent   bool      `gorm:"default:false" json:"is_percent"`          // 是否按百分比调整
	Priority    int       `gorm:"default:0" json:"priority"`                 // 优先级 (数字越大越高)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PricingRule) TableName() string {
	return "pricing_rules"
}
