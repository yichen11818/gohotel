package database

import (
	"fmt"
	"log"
	"time"

	"gohotel/internal/models"
	"gohotel/pkg/utils"
)

func seedDemoScenarioData() error {
	if err := seedDemoUsers(24); err != nil {
		return err
	}
	if err := seedDemoFacilities(12); err != nil {
		return err
	}
	if err := seedDemoPricingRules(4); err != nil {
		return err
	}
	if err := seedDemoBanners(6); err != nil {
		return err
	}
	if err := seedDemoNotices(4); err != nil {
		return err
	}
	if err := seedDemoBookings(10); err != nil {
		return err
	}
	if err := seedDemoWorkOrders(6, 8); err != nil {
		return err
	}
	if err := seedDemoBehaviors(120); err != nil {
		return err
	}

	return nil
}

func seedDemoUsers(minTotal int64) error {
	var userCount int64
	if err := DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return fmt.Errorf("统计用户数量失败: %w", err)
	}
	if userCount >= minTotal {
		log.Printf("⏩ 用户数据已达 %d 条，跳过补量", userCount)
		return nil
	}

	hashedPassword, err := utils.HashPassword("pass123")
	if err != nil {
		return fmt.Errorf("生成演示用户密码失败: %w", err)
	}

	profiles := []struct {
		Username   string
		Email      string
		RealName   string
		Phone      string
		Level      string
		Points     int
		Balance    float64
		TotalSpend float64
		Status     string
	}{
		{"demo_linhan", "linhan@gohotel.local", "林涵", "13900010001", "silver", 360, 128.50, 2680, "active"},
		{"demo_zhouwen", "zhouwen@gohotel.local", "周雯", "13900010002", "gold", 880, 560.00, 8560, "active"},
		{"demo_wangke", "wangke@gohotel.local", "王珂", "13900010003", "normal", 120, 0, 980, "active"},
		{"demo_suyun", "suyun@gohotel.local", "苏芸", "13900010004", "platinum", 1680, 1380.00, 16880, "active"},
		{"demo_qiuping", "qiuping@gohotel.local", "邱萍", "13900010005", "silver", 420, 220.00, 3180, "active"},
		{"demo_yexing", "yexing@gohotel.local", "叶星", "13900010006", "normal", 90, 0, 680, "active"},
		{"demo_litao", "litao@gohotel.local", "李涛", "13900010007", "gold", 1020, 760.00, 9280, "active"},
		{"demo_mengjie", "mengjie@gohotel.local", "孟洁", "13900010008", "normal", 60, 0, 520, "blocked"},
	}

	created := 0
	for _, profile := range profiles {
		if userCount+int64(created) >= minTotal {
			break
		}

		var exists int64
		if err := DB.Model(&models.User{}).Where("username = ?", profile.Username).Count(&exists).Error; err != nil {
			return fmt.Errorf("检查演示用户 %s 是否存在失败: %w", profile.Username, err)
		}
		if exists > 0 {
			continue
		}

		phone := profile.Phone
		user := models.User{
			ID:         utils.JSONInt64(utils.GenID()),
			Username:   profile.Username,
			Email:      profile.Email,
			Password:   hashedPassword,
			Phone:      &phone,
			RealName:   profile.RealName,
			Role:       "user",
			Status:     profile.Status,
			Level:      profile.Level,
			Points:     profile.Points,
			Balance:    profile.Balance,
			TotalSpend: profile.TotalSpend,
			FirstLogin: false,
		}

		if err := DB.Create(&user).Error; err != nil {
			return fmt.Errorf("创建演示用户 %s 失败: %w", profile.Username, err)
		}
		created++
	}

	if created > 0 {
		log.Printf("✅ 已补充 %d 条演示用户数据", created)
	} else {
		log.Printf("⏩ 演示用户模板已存在，当前用户总数 %d", userCount)
	}
	return nil
}

func seedDemoFacilities(minTotal int64) error {
	var facilityCount int64
	if err := DB.Model(&models.Facility{}).Count(&facilityCount).Error; err != nil {
		return fmt.Errorf("统计设施数量失败: %w", err)
	}
	if facilityCount >= minTotal {
		log.Printf("⏩ 设施数据已达 %d 条，跳过补量", facilityCount)
		return nil
	}

	templates := []models.Facility{
		{Type: "reception", Floor: 1, Left: 40, Top: 40, Width: 140, Height: 80, Label: "前台接待"},
		{Type: "elevator", Floor: 1, Left: 220, Top: 40, Width: 80, Height: 80, Label: "电梯 A"},
		{Type: "corridor", Floor: 1, Left: 320, Top: 56, Width: 260, Height: 60, Label: "主走廊"},
		{Type: "bathroom", Floor: 1, Left: 610, Top: 38, Width: 80, Height: 80, Label: "公共卫生间"},
		{Type: "laundry", Floor: 2, Left: 80, Top: 60, Width: 100, Height: 100, Label: "布草洗衣房"},
		{Type: "storage", Floor: 2, Left: 220, Top: 60, Width: 80, Height: 80, Label: "备品储物间"},
		{Type: "stairs", Floor: 2, Left: 340, Top: 40, Width: 80, Height: 120, Label: "消防楼梯"},
		{Type: "rest", Floor: 3, Left: 120, Top: 60, Width: 160, Height: 120, Label: "员工休息区"},
		{Type: "microwave", Floor: 3, Left: 320, Top: 70, Width: 100, Height: 100, Label: "茶歇餐饮区"},
	}

	created := 0
	for _, facility := range templates {
		if facilityCount+int64(created) >= minTotal {
			break
		}

		var exists int64
		if err := DB.Model(&models.Facility{}).
			Where("floor = ? AND label = ?", facility.Floor, facility.Label).
			Count(&exists).Error; err != nil {
			return fmt.Errorf("检查设施 %s 是否存在失败: %w", facility.Label, err)
		}
		if exists > 0 {
			continue
		}

		if err := DB.Create(&facility).Error; err != nil {
			return fmt.Errorf("创建设施 %s 失败: %w", facility.Label, err)
		}
		created++
	}

	if created > 0 {
		log.Printf("✅ 已补充 %d 条设施数据", created)
	} else {
		log.Printf("⏩ 设施模板已存在，当前设施总数 %d", facilityCount)
	}
	return nil
}

func seedDemoPricingRules(minTotal int64) error {
	var ruleCount int64
	if err := DB.Model(&models.PricingRule{}).Count(&ruleCount).Error; err != nil {
		return fmt.Errorf("统计定价规则数量失败: %w", err)
	}
	if ruleCount >= minTotal {
		log.Printf("⏩ 定价规则已达 %d 条，跳过补量", ruleCount)
		return nil
	}

	today := time.Now()
	baseDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	templates := []models.PricingRule{
		{Name: "周末弹性加价", Type: "weekend", RoomType: "", StartDate: baseDate, EndDate: baseDate.AddDate(0, 1, 0), Adjustment: 68, IsPercent: false, Priority: 4},
		{Name: "豪华套房会员礼遇", Type: "special", RoomType: "豪华套房", StartDate: baseDate, EndDate: baseDate.AddDate(0, 0, 45), Adjustment: -88, IsPercent: false, Priority: 6},
		{Name: "总统套房高峰档期", Type: "holiday", RoomType: "总统套房", StartDate: baseDate.AddDate(0, 0, 7), EndDate: baseDate.AddDate(0, 0, 14), Adjustment: 12, IsPercent: true, Priority: 7},
		{Name: "标准间连住促销", Type: "special", RoomType: "标准间", StartDate: baseDate, EndDate: baseDate.AddDate(0, 0, 30), Adjustment: -10, IsPercent: true, Priority: 5},
	}

	created := 0
	for _, rule := range templates {
		if ruleCount+int64(created) >= minTotal {
			break
		}

		var exists int64
		if err := DB.Model(&models.PricingRule{}).Where("name = ?", rule.Name).Count(&exists).Error; err != nil {
			return fmt.Errorf("检查定价规则 %s 是否存在失败: %w", rule.Name, err)
		}
		if exists > 0 {
			continue
		}

		if err := DB.Create(&rule).Error; err != nil {
			return fmt.Errorf("创建定价规则 %s 失败: %w", rule.Name, err)
		}
		created++
	}

	if created > 0 {
		log.Printf("✅ 已补充 %d 条定价规则", created)
	} else {
		log.Printf("⏩ 定价规则模板已存在，当前总数 %d", ruleCount)
	}
	return nil
}

func seedDemoBanners(minTotal int64) error {
	var bannerCount int64
	if err := DB.Model(&models.Banner{}).Count(&bannerCount).Error; err != nil {
		return fmt.Errorf("统计活动横幅数量失败: %w", err)
	}
	if bannerCount >= minTotal {
		log.Printf("⏩ 活动横幅已达 %d 条，跳过补量", bannerCount)
		return nil
	}

	subtitleA := "连住两晚享早餐与延迟退房"
	subtitleB := "商务住客可享会议室时段折扣"
	subtitleC := "家庭房套餐含儿童洗漱礼包"
	subtitleD := "金卡及以上会员积分双倍累计"
	subtitleE := "前台扫码即可查看入住须知"
	templates := []models.Banner{
		{Title: "春季住店礼遇", Subtitle: &subtitleA, ImageURL: "https://placehold.co/1200x420/png?text=GoHotel+Spring+Stay", Status: "active", Sort: 0},
		{Title: "商务差旅连住优惠", Subtitle: &subtitleB, ImageURL: "https://placehold.co/1200x420/png?text=GoHotel+Business+Stay", Status: "active", Sort: 0},
		{Title: "周末亲子套房推荐", Subtitle: &subtitleC, ImageURL: "https://placehold.co/1200x420/png?text=GoHotel+Family+Weekend", Status: "active", Sort: 0},
		{Title: "会员积分翻倍日", Subtitle: &subtitleD, ImageURL: "https://placehold.co/1200x420/png?text=GoHotel+Member+Day", Status: "active", Sort: 1},
		{Title: "自助入住服务指引", Subtitle: &subtitleE, ImageURL: "https://placehold.co/1200x420/png?text=GoHotel+Self+Checkin", Status: "active", Sort: 1},
	}

	created := 0
	for _, banner := range templates {
		if bannerCount+int64(created) >= minTotal {
			break
		}

		var exists int64
		if err := DB.Model(&models.Banner{}).Where("title = ?", banner.Title).Count(&exists).Error; err != nil {
			return fmt.Errorf("检查活动横幅 %s 是否存在失败: %w", banner.Title, err)
		}
		if exists > 0 {
			continue
		}

		if err := DB.Create(&banner).Error; err != nil {
			return fmt.Errorf("创建活动横幅 %s 失败: %w", banner.Title, err)
		}
		created++
	}

	if created > 0 {
		log.Printf("✅ 已补充 %d 条活动横幅", created)
	} else {
		log.Printf("⏩ 活动横幅模板已存在，当前总数 %d", bannerCount)
	}
	return nil
}

func seedDemoNotices(minTotal int64) error {
	var noticeCount int64
	if err := DB.Model(&models.Notice{}).Count(&noticeCount).Error; err != nil {
		return fmt.Errorf("统计公告数量失败: %w", err)
	}
	if noticeCount >= minTotal {
		log.Printf("⏩ 公告数据已达 %d 条，跳过补量", noticeCount)
		return nil
	}

	templates := []models.Notice{
		{Title: "入住时间调整提醒", Status: "active", Sort: 0},
		{Title: "发票开具说明", Status: "active", Sort: 1},
		{Title: "周末停车指引", Status: "active", Sort: 2},
		{Title: "会员早餐权益更新", Status: "active", Sort: 3},
	}

	created := 0
	for _, notice := range templates {
		if noticeCount+int64(created) >= minTotal {
			break
		}

		var exists int64
		if err := DB.Model(&models.Notice{}).Where("title = ?", notice.Title).Count(&exists).Error; err != nil {
			return fmt.Errorf("检查公告 %s 是否存在失败: %w", notice.Title, err)
		}
		if exists > 0 {
			continue
		}

		if err := DB.Create(&notice).Error; err != nil {
			return fmt.Errorf("创建公告 %s 失败: %w", notice.Title, err)
		}
		created++
	}

	if created > 0 {
		log.Printf("✅ 已补充 %d 条公告数据", created)
	} else {
		log.Printf("⏩ 公告模板已存在，当前总数 %d", noticeCount)
	}
	return nil
}

func seedDemoBookings(minTotal int64) error {
	var bookingCount int64
	if err := DB.Model(&models.Booking{}).Count(&bookingCount).Error; err != nil {
		return fmt.Errorf("统计订单数量失败: %w", err)
	}
	if bookingCount >= minTotal {
		log.Printf("⏩ 订单数据已达 %d 条，跳过补量", bookingCount)
		return nil
	}

	var users []models.User
	if err := DB.Where("role = ?", "user").Order("id ASC").Find(&users).Error; err != nil {
		return fmt.Errorf("查询演示用户失败: %w", err)
	}
	var rooms []models.Room
	if err := DB.Order("id ASC").Find(&rooms).Error; err != nil {
		return fmt.Errorf("查询演示房间失败: %w", err)
	}
	if len(users) == 0 || len(rooms) == 0 {
		log.Println("⏩ 缺少用户或房间数据，跳过订单补量")
		return nil
	}

	templates := []struct {
		Status         string
		PaymentStatus  string
		PaymentMethod  string
		CheckInOffset  int
		TotalDays      int
		Deposit        float64
		ExtraCharges   float64
		SpecialRequest string
		CancelReason   string
	}{
		{"pending", "unpaid", "", 2, 1, 0, 0, "需要安静房间", ""},
		{"confirmed", "paid", "wechat", 4, 2, 100, 0, "优先高楼层", ""},
		{"confirmed", "paid", "alipay", 7, 3, 100, 20, "商务连住，开具发票", ""},
		{"checkout", "paid", "card", -6, 2, 100, 35, "已完成离店结算", ""},
		{"cancelled", "refunded", "wechat", -3, 2, 0, 0, "", "客人临时改期"},
		{"confirmed", "paid", "alipay", 10, 1, 80, 0, "需要无烟房", ""},
		{"pending", "unpaid", "", 12, 2, 0, 0, "靠近电梯方便搬运行李", ""},
		{"checkout", "paid", "wechat", -10, 1, 50, 15, "已完成快速退房", ""},
		{"confirmed", "paid", "card", 1, 2, 120, 0, "带小朋友入住，请准备儿童拖鞋", ""},
		{"cancelled", "refunded", "alipay", 6, 1, 0, 0, "", "会议取消，申请退款"},
	}

	today := time.Now()
	created := 0
	for i := 0; bookingCount+int64(created) < minTotal && i < len(templates); i++ {
		template := templates[i]
		user := users[i%len(users)]
		room := rooms[i%len(rooms)]
		checkIn := time.Date(today.Year(), today.Month(), today.Day(), 15, 0, 0, 0, today.Location()).
			AddDate(0, 0, template.CheckInOffset)
		checkOut := checkIn.AddDate(0, 0, template.TotalDays)
		totalPrice := room.Price * float64(template.TotalDays)
		actualPrice := totalPrice + template.ExtraCharges

		guestName := user.RealName
		if guestName == "" {
			guestName = user.Username
		}
		guestPhone := fmt.Sprintf("138%08d", i+1001)
		if user.Phone != nil && *user.Phone != "" {
			guestPhone = *user.Phone
		}

		booking := models.Booking{
			ID:             utils.JSONInt64(utils.GenID()),
			BookingNumber:  utils.JSONInt64(utils.GenID()),
			UserID:         user.ID,
			RoomID:         int64(room.ID),
			CheckIn:        checkIn,
			CheckOut:       checkOut,
			TotalDays:      template.TotalDays,
			TotalPrice:     totalPrice,
			Deposit:        template.Deposit,
			ActualPrice:    actualPrice,
			ExtraCharges:   template.ExtraCharges,
			GuestName:      guestName,
			GuestPhone:     guestPhone,
			GuestIDCard:    fmt.Sprintf("42010219900101%04d", i+1001),
			ExtraGuests:    "[]",
			SpecialRequest: template.SpecialRequest,
			Status:         template.Status,
			PaymentStatus:  template.PaymentStatus,
			PaymentMethod:  template.PaymentMethod,
			CancelReason:   template.CancelReason,
			CreatedAt:      checkIn.Add(-48 * time.Hour),
			UpdatedAt:      checkIn.Add(-12 * time.Hour),
		}

		if err := DB.Create(&booking).Error; err != nil {
			return fmt.Errorf("创建演示订单失败: %w", err)
		}
		created++
	}

	if created > 0 {
		log.Printf("✅ 已补充 %d 条订单数据", created)
	} else {
		log.Printf("⏩ 订单模板已存在，当前总数 %d", bookingCount)
	}
	return nil
}

func seedDemoWorkOrders(minRepairs, minCleanings int64) error {
	var rooms []models.Room
	if err := DB.Order("id ASC").Find(&rooms).Error; err != nil {
		return fmt.Errorf("查询房间数据失败: %w", err)
	}
	if len(rooms) == 0 {
		log.Println("⏩ 没有房间数据，跳过工单补量")
		return nil
	}

	var admins []models.User
	if err := DB.Where("role = ? AND status = ?", "admin", "active").Order("id ASC").Find(&admins).Error; err != nil {
		return fmt.Errorf("查询管理员失败: %w", err)
	}
	var reporterID int64
	var workerID int64
	if len(admins) > 0 {
		reporterID = admins[0].ID.Int64()
		workerID = admins[len(admins)-1].ID.Int64()
	}

	var maintenanceCount int64
	if err := DB.Model(&models.Maintenance{}).Count(&maintenanceCount).Error; err != nil {
		return fmt.Errorf("统计维修工单数量失败: %w", err)
	}
	if maintenanceCount < minRepairs {
		templates := []models.Maintenance{
			{RoomID: int64(rooms[0%len(rooms)].ID), Type: "plumbing", Description: "淋浴热水忽冷忽热，需检查混水阀", Status: "pending", ReporterID: reporterID, CreatedAt: time.Now().Add(-10 * time.Hour)},
			{RoomID: int64(rooms[1%len(rooms)].ID), Type: "electrical", Description: "床头阅读灯闪烁，疑似接触不良", Status: "in_progress", ReporterID: reporterID, WorkerID: workerID, CreatedAt: time.Now().Add(-18 * time.Hour), UpdatedAt: time.Now().Add(-6 * time.Hour)},
			{RoomID: int64(rooms[2%len(rooms)].ID), Type: "furniture", Description: "行李架松动，需要重新加固", Status: "completed", ReporterID: reporterID, WorkerID: workerID, Remark: "已更换固定件并完成承重测试", CreatedAt: time.Now().Add(-30 * time.Hour), UpdatedAt: time.Now().Add(-20 * time.Hour)},
			{RoomID: int64(rooms[3%len(rooms)].ID), Type: "appliance", Description: "空调出风量偏小，客诉夜间闷热", Status: "pending", ReporterID: reporterID, CreatedAt: time.Now().Add(-5 * time.Hour)},
			{RoomID: int64(rooms[4%len(rooms)].ID), Type: "other", Description: "门锁识别偶发失败，建议检修", Status: "completed", ReporterID: reporterID, WorkerID: workerID, Remark: "已重置门锁模块并测试通过", CreatedAt: time.Now().Add(-48 * time.Hour), UpdatedAt: time.Now().Add(-36 * time.Hour)},
			{RoomID: int64(rooms[5%len(rooms)].ID), Type: "electrical", Description: "电视接口松动，信号源切换异常", Status: "pending", ReporterID: reporterID, CreatedAt: time.Now().Add(-2 * time.Hour)},
		}

		created := 0
		for _, maintenance := range templates {
			if maintenanceCount+int64(created) >= minRepairs {
				break
			}
			if err := DB.Create(&maintenance).Error; err != nil {
				return fmt.Errorf("创建维修工单失败: %w", err)
			}
			created++
		}
		if created > 0 {
			log.Printf("✅ 已补充 %d 条维修工单", created)
		}
	} else {
		log.Printf("⏩ 维修工单已达 %d 条，跳过补量", maintenanceCount)
	}

	var housekeepingCount int64
	if err := DB.Model(&models.Housekeeping{}).Count(&housekeepingCount).Error; err != nil {
		return fmt.Errorf("统计清洁工单数量失败: %w", err)
	}
	if housekeepingCount < minCleanings {
		now := time.Now()
		startA := now.Add(-4 * time.Hour)
		startB := now.Add(-8 * time.Hour)
		endB := now.Add(-6 * time.Hour)
		startC := now.Add(-2 * time.Hour)
		templates := []models.Housekeeping{
			{RoomID: int64(rooms[0%len(rooms)].ID), Type: "checkout", Status: "pending", CreatedAt: now.Add(-3 * time.Hour)},
			{RoomID: int64(rooms[1%len(rooms)].ID), Type: "daily", Status: "in_progress", StaffID: workerID, StartTime: &startA, CreatedAt: now.Add(-5 * time.Hour)},
			{RoomID: int64(rooms[2%len(rooms)].ID), Type: "deep", Status: "completed", StaffID: workerID, StartTime: &startB, EndTime: &endB, CreatedAt: now.Add(-10 * time.Hour), UpdatedAt: now.Add(-6 * time.Hour)},
			{RoomID: int64(rooms[3%len(rooms)].ID), Type: "checkout", Status: "pending", CreatedAt: now.Add(-90 * time.Minute)},
			{RoomID: int64(rooms[4%len(rooms)].ID), Type: "daily", Status: "completed", StaffID: workerID, StartTime: &startC, EndTime: &now, CreatedAt: now.Add(-3 * time.Hour), UpdatedAt: now},
			{RoomID: int64(rooms[5%len(rooms)].ID), Type: "daily", Status: "pending", CreatedAt: now.Add(-70 * time.Minute)},
			{RoomID: int64(rooms[6%len(rooms)].ID), Type: "deep", Status: "in_progress", StaffID: workerID, StartTime: &startC, CreatedAt: now.Add(-100 * time.Minute)},
			{RoomID: int64(rooms[7%len(rooms)].ID), Type: "checkout", Status: "pending", CreatedAt: now.Add(-40 * time.Minute)},
		}

		created := 0
		for _, housekeeping := range templates {
			if housekeepingCount+int64(created) >= minCleanings {
				break
			}
			if err := DB.Create(&housekeeping).Error; err != nil {
				return fmt.Errorf("创建清洁工单失败: %w", err)
			}
			created++
		}
		if created > 0 {
			log.Printf("✅ 已补充 %d 条清洁工单", created)
		}
	} else {
		log.Printf("⏩ 清洁工单已达 %d 条，跳过补量", housekeepingCount)
	}

	return nil
}

func seedDemoBehaviors(minTotal int64) error {
	var behaviorCount int64
	if err := DB.Model(&models.UserBehavior{}).Count(&behaviorCount).Error; err != nil {
		return fmt.Errorf("统计用户行为数量失败: %w", err)
	}
	if behaviorCount >= minTotal {
		log.Printf("⏩ 用户行为数据已达 %d 条，跳过补量", behaviorCount)
		return nil
	}

	var users []models.User
	if err := DB.Where("role = ? AND status = ?", "user", "active").Order("id ASC").Find(&users).Error; err != nil {
		return fmt.Errorf("查询行为用户失败: %w", err)
	}
	var rooms []models.Room
	if err := DB.Order("id ASC").Find(&rooms).Error; err != nil {
		return fmt.Errorf("查询行为房间失败: %w", err)
	}
	if len(users) == 0 || len(rooms) == 0 {
		log.Println("⏩ 缺少用户或房间数据，跳过行为补量")
		return nil
	}

	created := int64(0)
	baseTime := time.Now()
	for round := 0; behaviorCount+created < minTotal; round++ {
		for userIndex, user := range users {
			if behaviorCount+created >= minTotal {
				break
			}

			primaryRoom := rooms[(userIndex+round)%len(rooms)]
			secondaryRoom := rooms[(userIndex+round+1)%len(rooms)]
			events := []struct {
				roomID       int64
				behaviorType string
				source       string
				weight       float64
			}{
				{int64(primaryRoom.ID), "view_room", "h5_home", 1.0},
				{int64(primaryRoom.ID), "view_detail", "h5_detail", 1.8},
				{int64(secondaryRoom.ID), "click_recommendation", "h5_recommendation", 2.2},
			}
			if (userIndex+round)%2 == 0 {
				events = append(events, struct {
					roomID       int64
					behaviorType string
					source       string
					weight       float64
				}{int64(primaryRoom.ID), "book_intent", "h5_booking", 3.6})
			}

			for eventIndex, event := range events {
				if behaviorCount+created >= minTotal {
					break
				}

				behavior := models.UserBehavior{
					ID:           utils.JSONInt64(utils.GenID()),
					UserID:       user.ID,
					RoomID:       event.roomID,
					BehaviorType: event.behaviorType,
					Source:       event.source,
					Weight:       event.weight,
					MetadataJSON: fmt.Sprintf(`{"seed":"demo","round":%d,"event_index":%d}`, round, eventIndex),
					CreatedAt:    baseTime.Add(-time.Duration(userIndex+eventIndex+round) * time.Hour),
					UpdatedAt:    baseTime.Add(-time.Duration(userIndex+eventIndex+round) * time.Hour),
				}

				if err := DB.Create(&behavior).Error; err != nil {
					return fmt.Errorf("创建用户行为失败: %w", err)
				}
				created++
			}
		}
	}

	if created > 0 {
		log.Printf("✅ 已补充 %d 条用户行为数据", created)
	}
	return nil
}
