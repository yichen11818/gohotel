package models

import (
	"time"
)

// Room 房间模型
// 对应数据库中的 rooms 表
type Room struct {
	ID            uint      `gorm:"primaryKey" json:"id"`                             // 主键
	RoomNumber    string    `gorm:"unique;not null;size:20;index" json:"room_number"` // 房间号（唯一，有索引）
	RoomType      string    `gorm:"not null;size:50;index" json:"room_type"`          // 房间类型（有索引）
	Floor         int       `gorm:"not null" json:"floor"`                            // 楼层
	Price         float64   `gorm:"not null;type:decimal(10,2)" json:"price"`         // 价格（每晚）
	OriginalPrice float64   `gorm:"type:decimal(10,2)" json:"original_price"`         // 原价
	Capacity      int       `gorm:"not null" json:"capacity"`                         // 可住人数
	Area          float64   `gorm:"type:decimal(10,2)" json:"area"`                   // 面积（平方米）
	BedType       string    `gorm:"size:50" json:"bed_type"`                          // 床型：单人床、双人床、大床
	Description   string    `gorm:"type:text" json:"description"`                     // 房间描述
	Facilities    string    `gorm:"type:text" json:"facilities"`                      // 设施（JSON 字符串）
	Images        string    `gorm:"type:text" json:"images"`                          // 图片 URL（JSON 数组）
	Status        string    `gorm:"default:'available';size:20;index" json:"status"`  // 状态：available, occupied, maintenance
	CreatedAt     time.Time `json:"created_at"`                                       // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`                                       // 更新时间
}

// TableName 指定表名
func (Room) TableName() string {
	return "rooms"
}

// IsAvailable 判断房间是否可用
func (r *Room) IsAvailable() bool {
	return r.Status == "available"
}

// GetDiscountRate 获取折扣率
func (r *Room) GetDiscountRate() float64 {
	if r.OriginalPrice == 0 {
		return 0
	}
	return (r.OriginalPrice - r.Price) / r.OriginalPrice * 100
}

