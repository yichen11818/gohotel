package database

import (
	"fmt"
	"log"

	"gohotel/internal/models"
)

// AutoMigrate 自动迁移数据库
// 这个函数会根据模型自动创建/更新表结构
func AutoMigrate() error {
	log.Println("🔄 开始数据库迁移...")

	// AutoMigrate 会：
	// 1. 创建不存在的表
	// 2. 添加缺失的列
	// 3. 添加缺失的索引
	// 注意：不会删除已存在的列（为了安全）
	err := DB.AutoMigrate(
		&models.User{},
		&models.Hotel{},
		&models.HotelSettings{},
		&models.Room{},
		&models.Booking{},
		&models.Facility{},
		&models.Log{},
		&models.Banner{},
		&models.Notice{},
	)

	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	log.Println("✅ 数据库迁移完成！")
	return nil
}

// SeedData 插入测试数据（可选）
// 在开发环境中，可以用这个函数插入一些示例数据
func SeedData() error {
	log.Println("🌱 开始插入测试数据...")

	// 检查是否已有数据
	var roomCount int64
	DB.Model(&models.Room{}).Count(&roomCount)
	if roomCount > 0 {
		log.Println("⏩ 已存在房间数据，跳过插入")
		return nil
	}

	// 插入示例房间
	rooms := []models.Room{
		{
			RoomNumber:    "101",
			RoomType:      "标准间",
			Floor:         1,
			Price:         200.00,
			OriginalPrice: 280.00,
			Capacity:      2,
			Area:          25.0,
			BedType:       "双床",
			Description:   "舒适的标准双人间，配有独立卫浴和空调",
			Facilities:    `["WiFi", "空调", "电视", "热水器"]`,
			Status:        "available",
		},
		{
			RoomNumber:    "201",
			RoomType:      "豪华套房",
			Floor:         2,
			Price:         500.00,
			OriginalPrice: 680.00,
			Capacity:      4,
			Area:          45.0,
			BedType:       "大床+沙发床",
			Description:   "宽敞的豪华套房，带客厅和阳台，视野开阔",
			Facilities:    `["WiFi", "空调", "电视", "热水器", "浴缸", "阳台"]`,
			Status:        "available",
		},
		{
			RoomNumber:    "301",
			RoomType:      "总统套房",
			Floor:         3,
			Price:         1000.00,
			OriginalPrice: 1500.00,
			Capacity:      6,
			Area:          80.0,
			BedType:       "特大床",
			Description:   "顶级总统套房，配有私人管家服务和独立会客厅",
			Facilities:    `["WiFi", "空调", "电视", "热水器", "浴缸", "阳台", "音响", "投影仪"]`,
			Status:        "available",
		},
	}

	if err := DB.Create(&rooms).Error; err != nil {
		return fmt.Errorf("插入房间数据失败: %w", err)
	}

	log.Printf("✅ 成功插入 %d 条房间数据", len(rooms))
	return nil
}
