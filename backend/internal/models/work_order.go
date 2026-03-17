package models

import (
	"time"
)

// Maintenance 维修/报修模型
type Maintenance struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoomID      int64     `gorm:"not null;index" json:"room_id"`
	Type        string    `gorm:"not null;size:50" json:"type"`             // 报修类型：plumbing, electrical, furniture, etc.
	Description string    `gorm:"type:text" json:"description"`             // 问题描述
	Status      string    `gorm:"default:'pending';size:20" json:"status"` // 状态：pending, in_progress, completed, cancelled
	ReporterID  int64     `gorm:"index" json:"reporter_id"`                // 报修人 (员工 ID)
	WorkerID    int64     `gorm:"index" json:"worker_id"`                  // 维修人 (员工 ID)
	Remark      string    `gorm:"type:text" json:"remark"`                 // 维修备注
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Room Room `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}

func (Maintenance) TableName() string {
	return "maintenances"
}

// Housekeeping 清洁任务模型
type Housekeeping struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoomID      int64     `gorm:"not null;index" json:"room_id"`
	Type        string    `gorm:"size:50" json:"type"`                      // 类型：daily(续住清), checkout(退房清), deep(深清)
	Status      string    `gorm:"default:'pending';size:20" json:"status"` // 状态：pending, in_progress, completed
	StaffID     int64     `gorm:"index" json:"staff_id"`                   // 清洁人员 ID
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Room Room `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}

func (Housekeeping) TableName() string {
	return "housekeepings"
}
