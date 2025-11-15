package models

import (
	"time"
)

// Booking 预订模型
// 对应数据库中的 bookings 表
type Booking struct {
	ID             uint      `gorm:"primaryKey" json:"id"`                                 // 主键
	BookingNumber  string    `gorm:"unique;not null;size:50;index" json:"booking_number"`  // 预订单号（唯一）
	UserID         uint      `gorm:"not null;index" json:"user_id"`                        // 用户 ID（有索引）
	RoomID         uint      `gorm:"not null;index" json:"room_id"`                        // 房间 ID（有索引）
	CheckIn        time.Time `gorm:"not null;index" json:"check_in"`                       // 入住日期（有索引）
	CheckOut       time.Time `gorm:"not null;index" json:"check_out"`                      // 退房日期（有索引）
	TotalDays      int       `gorm:"not null" json:"total_days"`                           // 总天数
	TotalPrice     float64   `gorm:"not null;type:decimal(10,2)" json:"total_price"`       // 总价
	GuestName      string    `gorm:"not null;size:50" json:"guest_name"`                   // 入住人姓名
	GuestPhone     string    `gorm:"not null;size:20" json:"guest_phone"`                  // 入住人电话
	GuestIDCard    string    `gorm:"size:50" json:"guest_id_card"`                         // 入住人身份证号
	SpecialRequest string    `gorm:"type:text" json:"special_request"`                     // 特殊要求
	Status         string    `gorm:"default:'pending';size:20;index" json:"status"`        // 状态：pending, confirmed, checkin, checkout, cancelled
	PaymentStatus  string    `gorm:"default:'unpaid';size:20;index" json:"payment_status"` // 支付状态：unpaid, paid, refunded
	PaymentMethod  string    `gorm:"size:50" json:"payment_method"`                        // 支付方式：wechat, alipay, card
	CancelReason   string    `gorm:"type:text" json:"cancel_reason"`                       // 取消原因
	CreatedAt      time.Time `json:"created_at"`                                           // 创建时间
	UpdatedAt      time.Time `json:"updated_at"`                                           // 更新时间

	// 关联查询（可选）
	// 当查询 Booking 时，可以同时加载 User 和 Room 的信息
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"` // 关联的用户
	Room Room `gorm:"foreignKey:RoomID" json:"room,omitempty"` // 关联的房间
}

// TableName 指定表名
func (Booking) TableName() string {
	return "bookings"
}

// IsPending 判断是否待确认
func (b *Booking) IsPending() bool {
	return b.Status == "pending"
}

// IsConfirmed 判断是否已确认
func (b *Booking) IsConfirmed() bool {
	return b.Status == "confirmed"
}

// IsCancelled 判断是否已取消
func (b *Booking) IsCancelled() bool {
	return b.Status == "cancelled"
}

// IsPaid 判断是否已支付
func (b *Booking) IsPaid() bool {
	return b.PaymentStatus == "paid"
}

// CanCancel 判断是否可以取消
// 只有待确认和已确认的订单可以取消
func (b *Booking) CanCancel() bool {
	return b.Status == "pending" || b.Status == "confirmed"
}

// CanCheckIn 判断是否可以入住
// 已确认且已支付的订单可以入住
func (b *Booking) CanCheckIn() bool {
	return b.Status == "confirmed" && b.IsPaid()
}
