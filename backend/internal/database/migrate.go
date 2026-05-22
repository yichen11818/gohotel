package database

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"gohotel/internal/models"
	"gohotel/pkg/utils"

	"gorm.io/gorm"
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
		&models.RoomCategory{},
		&models.Room{},
		&models.Booking{},
		&models.Facility{},
		&models.Log{},
		&models.Banner{},
		&models.Notice{},
		&models.RoomInventory{},
		&models.UserBehavior{},
		&models.Maintenance{},
		&models.Housekeeping{},
		&models.PricingRule{},
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

	hotelID, err := seedDefaultHotel()
	if err != nil {
		return err
	}

	if err := seedDefaultHotelSettings(hotelID); err != nil {
		return err
	}

	if err := seedDefaultAdmin(); err != nil {
		return err
	}

	roomTemplates := defaultSeedRooms()

	// 检查是否已有房间数据
	var roomCount int64
	DB.Model(&models.Room{}).Count(&roomCount)
	if roomCount > 0 {
		log.Println("⏩ 已存在房间数据，跳过插入")
		if err := backfillSeedRoomDisplayData(roomTemplates); err != nil {
			return err
		}
		if err := seedRoomCategoriesFromRooms(); err != nil {
			return err
		}
		if err := seedDefaultInventory(); err != nil {
			return err
		}
		return seedDemoScenarioData()
	}

	// 插入示例房间
	rooms := roomTemplates

	if err := DB.Create(&rooms).Error; err != nil {
		return fmt.Errorf("插入房间数据失败: %w", err)
	}

	log.Printf("✅ 成功插入 %d 条房间数据", len(rooms))
	if err := seedRoomCategoriesFromRooms(); err != nil {
		return err
	}
	if err := seedDefaultInventory(); err != nil {
		return err
	}
	return seedDemoScenarioData()
}

func seedRoomCategoriesFromRooms() error {
	var rooms []models.Room
	if err := DB.Order("id ASC").Find(&rooms).Error; err != nil {
		return fmt.Errorf("查询房间数据失败: %w", err)
	}
	if len(rooms) == 0 {
		return nil
	}

	for _, room := range rooms {
		var existing models.RoomCategory
		err := DB.Where("name = ?", room.RoomType).First(&existing).Error
		if err == nil {
			updates := map[string]interface{}{}
			if strings.TrimSpace(existing.Description) == "" && strings.TrimSpace(room.Description) != "" {
				updates["description"] = room.Description
			}

			roomFacilities := normalizeSeedJSONArray(room.Facilities)
			if normalizeSeedJSONArray(existing.Facilities) == "[]" && roomFacilities != "[]" {
				updates["facilities"] = roomFacilities
			}

			roomImages := normalizeSeedJSONArray(room.Images)
			if normalizeSeedJSONArray(existing.Images) == "[]" && roomImages != "[]" {
				updates["images"] = roomImages
			}

			if len(updates) > 0 {
				if err := DB.Model(&existing).Updates(updates).Error; err != nil {
					return fmt.Errorf("补充房型分类 %s 失败: %w", room.RoomType, err)
				}
			}
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("检查房型分类 %s 失败: %w", room.RoomType, err)
		}

		category := models.RoomCategory{
			Name:        room.RoomType,
			Description: room.Description,
			Facilities:  normalizeSeedJSONArray(room.Facilities),
			Images:      normalizeSeedJSONArray(room.Images),
		}
		if err := DB.Create(&category).Error; err != nil {
			return fmt.Errorf("创建房型分类 %s 失败: %w", room.RoomType, err)
		}
	}

	return nil
}

func defaultSeedRooms() []models.Room {
	return []models.Room{
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
			Description:   "宽敞的豪华套房配备独立会客区与观景阳台，兼顾商务会客和家庭入住，采光通透，适合品质度假。",
			Facilities:    `["WiFi", "空调", "电视", "热水器", "浴缸", "阳台", "独立客厅", "迷你吧", "咖啡机"]`,
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
}

func backfillSeedRoomDisplayData(templates []models.Room) error {
	updatedFields := 0

	for _, template := range templates {
		if strings.TrimSpace(template.RoomType) == "" {
			continue
		}

		if strings.TrimSpace(template.Description) != "" {
			result := DB.Model(&models.Room{}).
				Where("room_type = ?", template.RoomType).
				Where("description IS NULL OR TRIM(description) = ''").
				Update("description", template.Description)
			if result.Error != nil {
				return fmt.Errorf("补充房间 %s 描述失败: %w", template.RoomType, result.Error)
			}
			updatedFields += int(result.RowsAffected)
		}

		facilities := normalizeSeedJSONArray(template.Facilities)
		if facilities != "[]" {
			result := DB.Model(&models.Room{}).
				Where("room_type = ?", template.RoomType).
				Where("facilities IS NULL OR TRIM(facilities) = '' OR TRIM(facilities) = '[]'").
				Update("facilities", facilities)
			if result.Error != nil {
				return fmt.Errorf("补充房间 %s 设施失败: %w", template.RoomType, result.Error)
			}
			updatedFields += int(result.RowsAffected)
		}

		var category models.RoomCategory
		err := DB.Where("name = ?", template.RoomType).First(&category).Error
		if err == gorm.ErrRecordNotFound {
			continue
		}
		if err != nil {
			return fmt.Errorf("查询房型分类 %s 失败: %w", template.RoomType, err)
		}

		updates := map[string]interface{}{}
		if strings.TrimSpace(category.Description) == "" && strings.TrimSpace(template.Description) != "" {
			updates["description"] = template.Description
		}
		if normalizeSeedJSONArray(category.Facilities) == "[]" && facilities != "[]" {
			updates["facilities"] = facilities
		}

		if len(updates) == 0 {
			continue
		}

		result := DB.Model(&category).Updates(updates)
		if result.Error != nil {
			return fmt.Errorf("补充房型分类 %s 展示信息失败: %w", template.RoomType, result.Error)
		}
		updatedFields += int(result.RowsAffected)
	}

	if updatedFields > 0 {
		log.Printf("✅ 已补充 %d 个房间/房型展示字段", updatedFields)
	}

	return nil
}

func normalizeSeedJSONArray(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return "[]"
	}

	var parsed []string
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return "[]"
	}

	normalized, err := json.Marshal(parsed)
	if err != nil {
		return "[]"
	}

	return string(normalized)
}

func seedDefaultAdmin() error {
	var adminCount int64
	if err := DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount).Error; err != nil {
		return fmt.Errorf("检查管理员数据失败: %w", err)
	}
	if adminCount > 0 {
		log.Println("⏩ 已存在管理员账号，跳过初始化")
		return nil
	}

	hashedPassword, err := utils.HashPassword("Admin@123456")
	if err != nil {
		return fmt.Errorf("初始化管理员密码失败: %w", err)
	}

	adminUser := models.User{
		ID:         utils.JSONInt64(utils.GenID()),
		Username:   "admin",
		Email:      "admin@gohotel.local",
		Password:   hashedPassword,
		RealName:   "系统管理员",
		Role:       "admin",
		Status:     "active",
		FirstLogin: true,
	}

	if err := DB.Create(&adminUser).Error; err != nil {
		return fmt.Errorf("插入默认管理员失败: %w", err)
	}

	log.Println("✅ 已初始化默认管理员账号：admin / Admin@123456")
	return nil
}

func seedDefaultHotel() (int64, error) {
	var hotel models.Hotel
	err := DB.Order("id ASC").First(&hotel).Error
	if err == nil {
		log.Printf("⏩ 已存在酒店数据，使用酒店 #%d", hotel.ID)
		return hotel.ID, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, fmt.Errorf("检查酒店数据失败: %w", err)
	}

	hotel = models.Hotel{
		Name:   "七天酒店",
		Status: "active",
	}
	if err := DB.Create(&hotel).Error; err != nil {
		return 0, fmt.Errorf("插入默认酒店失败: %w", err)
	}

	log.Printf("✅ 已初始化默认酒店：%s (#%d)", hotel.Name, hotel.ID)
	return hotel.ID, nil
}

func seedDefaultHotelSettings(hotelID int64) error {
	var settings models.HotelSettings
	err := DB.Where("hotel_id = ?", hotelID).First(&settings).Error
	if err == nil {
		log.Printf("⏩ 酒店 #%d 已存在酒店设置，跳过初始化", hotelID)
		return nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("检查酒店设置失败: %w", err)
	}

	defaultSettings := map[string]interface{}{
		"hotel_profile": map[string]interface{}{
			"name":              "七天酒店",
			"brand_name":        "GoHotel",
			"address":           "欢迎入住七天酒店",
			"front_desk_phone":  "400-800-1234",
			"service_time_text": "24小时服务",
			"intro_html":        "<p>欢迎入住七天酒店，享受舒适便捷的入住体验。</p>",
			"cover_images":      []string{},
			"geo": map[string]interface{}{
				"lat": 30.56,
				"lng": 114.28,
			},
		},
		"booking_rules": map[string]interface{}{
			"check_in_time":          "14:00",
			"check_out_time":         "12:00",
			"min_advance_days":       0,
			"max_advance_days":       30,
			"min_nights":             1,
			"max_nights":             30,
			"allow_same_day_booking": true,
		},
	}

	payload, err := json.Marshal(defaultSettings)
	if err != nil {
		return fmt.Errorf("序列化默认酒店设置失败: %w", err)
	}

	settings = models.HotelSettings{
		HotelID:      hotelID,
		SettingsJSON: string(payload),
	}
	if err := DB.Create(&settings).Error; err != nil {
		return fmt.Errorf("插入默认酒店设置失败: %w", err)
	}

	log.Printf("✅ 已初始化酒店 #%d 的默认设置", hotelID)
	return nil
}

func seedDefaultInventory() error {
	var rooms []models.Room
	if err := DB.Find(&rooms).Error; err != nil {
		return fmt.Errorf("查询房间数据失败: %w", err)
	}
	if len(rooms) == 0 {
		log.Println("⏩ 暂无房间数据，跳过库存初始化")
		return nil
	}

	type inventorySeedMeta struct {
		totalCount int
		price      float64
	}

	roomTypeMeta := make(map[string]inventorySeedMeta)
	for _, room := range rooms {
		meta := roomTypeMeta[room.RoomType]
		meta.totalCount++
		if meta.price == 0 {
			meta.price = room.Price
		}
		roomTypeMeta[room.RoomType] = meta
	}

	today := time.Now()
	startDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	const days = 30

	for roomType, meta := range roomTypeMeta {
		createdCount := 0
		for i := 0; i < days; i++ {
			date := startDate.AddDate(0, 0, i)

			var existingCount int64
			nextDate := date.AddDate(0, 0, 1)
			if err := DB.Model(&models.RoomInventory{}).
				Where("room_type = ? AND date >= ? AND date < ?", roomType, date.Format("2006-01-02"), nextDate.Format("2006-01-02")).
				Count(&existingCount).Error; err != nil {
				return fmt.Errorf("检查房型 %s 在 %s 的库存失败: %w", roomType, date.Format("2006-01-02"), err)
			}
			if existingCount > 0 {
				continue
			}

			inventory := models.RoomInventory{
				RoomType:   roomType,
				Date:       date,
				TotalCount: meta.totalCount,
				Price:      meta.price,
			}
			if err := DB.Create(&inventory).Error; err != nil {
				return fmt.Errorf("初始化房型 %s 在 %s 的库存失败: %w", roomType, date.Format("2006-01-02"), err)
			}

			createdCount++
		}

		if createdCount == 0 {
			log.Printf("⏩ 房型 %s 的默认库存已存在，跳过初始化", roomType)
			continue
		}

		log.Printf("✅ 已为房型 %s 初始化未来 %d 天库存，房量=%d", roomType, createdCount, meta.totalCount)
	}

	return nil
}
