package models

import (
	"time"
)

// RoomInventory 房量库存模型
// 对应数据库中的 room_inventories 表
type RoomInventory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RoomType   string    `gorm:"not null;size:50;index" json:"room_type"` // 房型
	Date       time.Time `gorm:"not null;index" json:"date"`              // 日期
	TotalCount int       `gorm:"not null" json:"total_count"`             // 总量
	BookedCount int      `gorm:"not null;default:0" json:"booked_count"`  // 已售数量
	Price      float64   `gorm:"not null;type:decimal(10,2)" json:"price"` // 当日价格
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (RoomInventory) TableName() string {
	return "room_inventories"
}
