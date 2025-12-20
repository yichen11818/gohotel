package main

import (
	"fmt"
	"log"
	"time"

	"gohotel/internal/config"
	"gohotel/internal/database"
	"gohotel/internal/handler"
	"gohotel/internal/middleware"
	"gohotel/internal/repository"
	"gohotel/internal/service"
	"gohotel/pkg/logger"
	"gohotel/pkg/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gohotel/docs" // 导入生成的 docs 包
)

// @title           酒店管理系统 API
// @version         1.0
// @description     这是一个酒店预订管理系统的 RESTful API 服务
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@gohotel.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      nas.yumi.chat:19999
// @BasePath

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 输入 "Bearer {token}" 格式的 JWT token

func main() {
	// 1. 加载配置
	fmt.Println("📖 正在加载配置...")
	if err := config.Load(); err != nil {
		log.Fatal("配置加载失败:", err)
	}
	fmt.Println("✅ 配置加载成功!")

	// 2. 初始化日志
	fmt.Println("📝 正在初始化日志...")
	if err := logger.Init(&logger.LogConfig{
		Level:      config.AppConfig.Log.Level,
		Filename:   config.AppConfig.Log.Filename,
		MaxSize:    config.AppConfig.Log.MaxSize,
		MaxBackups: config.AppConfig.Log.MaxBackups,
		MaxAge:     config.AppConfig.Log.MaxAge,
		Compress:   config.AppConfig.Log.Compress,
		Console:    config.AppConfig.Log.Console,
	}); err != nil {
		log.Fatal("日志初始化失败:", err)
	}
	defer logger.Sync()
	fmt.Println("✅ 日志初始化成功!")

	// 3. 连接数据库
	fmt.Println("🔌 正在连接数据库...")
	if err := database.InitMySQL(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer database.CloseDB()

	// 4. 自动迁移数据库表
	fmt.Println("🔄 正在执行数据库迁移...")
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 5. 插入测试数据（可选）
	if err := database.SeedData(); err != nil {
		log.Fatal("测试数据插入失败:", err)
	}

	// 6. 初始化雪花算法节点
	fmt.Println("❄️  正在初始化雪花算法节点...")
	// 节点ID可以从配置文件读取，这里暂时使用固定值 1
	// 如果是分布式部署，需要确保每个实例使用不同的节点ID（0-1023）
	if err := utils.InitSnowflake(1); err != nil {
		log.Fatal("雪花算法初始化失败:", err)
	}
	fmt.Println("✅ 雪花算法初始化成功!")

	// 6.1 初始化COS服务
	fmt.Println("☁️  正在初始化COS服务...")
	var cosService *service.CosService
	var cosErr error
	cosService, cosErr = service.NewCosService(&config.AppConfig.COS)
	if cosErr != nil {
		log.Printf("⚠️  COS服务初始化失败: %v，将无法使用图片上传功能", cosErr)
		cosService = nil
	} else {
		fmt.Println("✅ COS服务初始化成功!")
	}

	// 6.2 初始化时间轮
	fmt.Println("⏰ 正在初始化时间轮...")
	timeWheel := utils.NewMultiTimeWheel() // 使用多层时间轮（秒、分、时、天四层）

	// 设置持久化存储，将任务保存到data目录
	persistStore := utils.NewFilePersistStore("./data/timewheel_tasks.json")
	timeWheel.SetPersistStore(persistStore)

	timeWheel.Start()
	defer timeWheel.Stop()
	fmt.Println("✅ 时间轮初始化成功!")

	// 7. 初始化依赖注入
	// Repository 层
	userRepo := repository.NewUserRepository(database.DB)
	roomRepo := repository.NewRoomRepository(database.DB)
	bookingRepo := repository.NewBookingRepository(database.DB)
	logRepo := repository.NewLogRepository(database.DB)
	facilityRepo := repository.NewFacilityRepository(database.DB)
	bannerRepo := repository.NewBannerRepository(database.DB)

	// Service 层
	userService := service.NewUserService(userRepo)
	roomService := service.NewRoomService(roomRepo)
	bookingService := service.NewBookingService(bookingRepo, roomRepo, userRepo)
	logService := service.NewLogService(logRepo)
	facilityService := service.NewFacilityService(facilityRepo)
	bannerService := service.NewBannerService(bannerRepo, cosService, timeWheel)

	// 加载持久化的时间轮任务
	fmt.Println("📂 正在加载时间轮任务...")
	if err := timeWheel.LoadTasks(); err != nil {
		log.Printf("⚠️  时间轮任务加载失败: %v", err)
	} else {
		fmt.Println("✅ 时间轮任务加载成功!")
	}

	// 添加COS临时文件清理任务
	if cosService != nil {
		// 立即执行一次清理
		go func() {
			_, err := cosService.CleanupTempFiles(1 * time.Hour) // 清理1小时前的临时文件
			if err != nil {
				log.Printf("⚠️  初始清理COS临时文件失败: %v", err)
			}
		}()

		// 使用命名函数来实现递归调用
		var cleanupTask func()
		cleanupTask = func() {
			_, err := cosService.CleanupTempFiles(1 * time.Hour) // 清理1小时前的临时文件
			if err != nil {
				log.Printf("⚠️  清理COS临时文件失败: %v", err)
			}

			// 任务执行完成后，添加下一次任务
			nextExecTime := time.Now().Add(30 * time.Minute)
			timeWheel.AddTask(nextExecTime, cleanupTask, nil)
		}

		// 计算下一次执行时间
		nextExecTime := time.Now().Add(30 * time.Minute)

		// 添加任务到时间轮
		timeWheel.AddTask(nextExecTime, cleanupTask, nil)

		fmt.Println("✅ COS临时文件清理任务已添加，每30分钟执行一次")
	}

	// Handler 层
	userHandler := handler.NewUserHandler(userService)
	roomHandler := handler.NewRoomHandler(roomService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	logHandler := handler.NewLogHandler(logService)
	facilityHandler := handler.NewFacilityHandler(facilityService)
	bannerHandler := handler.NewBannerHandler(bannerService, cosService)
	cosHandler := handler.NewCosHandler(cosService)

	// 8. 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 9. 创建 Gin 引擎
	r := gin.New()

	// 10. 使用中间件
	r.Use(gin.Recovery())                // 恢复中间件（处理 panic）
	r.Use(middleware.CORSMiddleware())   // 跨域中间件
	r.Use(middleware.LoggerMiddleware()) // 日志中间件

	// 设置路由
	setupRoutes(r, userHandler, roomHandler, bookingHandler, logHandler, facilityHandler, bannerHandler, cosHandler)

	// 12. 启动服务器
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("🏨 酒店管理系统 API 服务器")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Printf("📍 服务器地址: http://%s\n", config.AppConfig.Server.Port)
	fmt.Printf("📝 运行模式: %s\n", config.AppConfig.Server.Mode)
	fmt.Printf("📚 Swagger 文档: http://%s/swagger/index.html\n", config.AppConfig.Server.Port)
	fmt.Println("═══════════════════════════════════════════════")

	if err := r.Run(config.AppConfig.Server.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

// setupRoutes 设置所有路由
func setupRoutes(r *gin.Engine, userHandler *handler.UserHandler, roomHandler *handler.RoomHandler, bookingHandler *handler.BookingHandler, logHandler *handler.LogHandler, facilityHandler *handler.FacilityHandler, bannerHandler *handler.BannerHandler, cosHandler *handler.CosHandler) {
	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	// @Summary 健康检查
	// @Description 检查服务器运行状态
	// @Tags 系统
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "酒店管理系统运行正常",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 认证路由（公开）
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// 房间路由（公开查询）
		rooms := api.Group("/rooms")
		{
			rooms.GET("", roomHandler.ListRooms)                     // 获取所有房间
			rooms.GET("/available", roomHandler.ListAvailableRooms)  // 获取可用房间
			rooms.GET("/floor/:floor", roomHandler.GetRoomByFloor)   // 根据楼层获取房间
			rooms.GET("/search/type", roomHandler.SearchRoomsByType) // 按房型搜索
			rooms.GET("/:id", roomHandler.GetRoomByID)               // 获取房间详情

			// 需要认证的房间管理路由（管理员）
			roomsAuth := rooms.Group("")
			roomsAuth.Use(middleware.AuthMiddleware())
			{
				roomsAuth.POST("", roomHandler.CreateRoom)             // 创建房间
				roomsAuth.POST("/batch", roomHandler.BatchCreateRooms) // 批量创建房间
				roomsAuth.POST("/:id", roomHandler.UpdateRoom)         // 更新房间
				roomsAuth.POST("/:id/delete", roomHandler.DeleteRoom)  // 删除房间
			}
		}
		// 活动横幅路由（公开查询）
		banners := api.Group("/banners")
		{
			banners.GET("/active", bannerHandler.GetActiveBanners) // 获取激活的活动横幅（前端展示用）
		}
		// 日志路由
		logs := api.Group("/logs")
		{
			logs.POST("/report", logHandler.Report) // 上报日志
			logs.GET("", logHandler.GetLogs)        // 获取日志列表
		}

		// 文件上传路由（需要认证，但不需要管理员权限）
		upload := api.Group("/upload")
		upload.Use(middleware.AuthMiddleware())
		{
			upload.POST("/image", cosHandler.UploadImage) // 通用图片上传接口
		}

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 用户路由
			users := authorized.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)       // 获取个人信息
				users.POST("/profile", userHandler.UpdateProfile)   // 更新个人信息
				users.POST("/password", userHandler.ChangePassword) // 修改密码
			}

			// 预订路由
			bookings := authorized.Group("/bookings")
			{
				bookings.POST("", bookingHandler.CreateBooking)            // 创建预订
				bookings.GET("/my", bookingHandler.GetMyBookings)          // 我的预订列表
				bookings.GET("/:id", bookingHandler.GetBookingByID)        // 获取预订详情
				bookings.POST("/:id/cancel", bookingHandler.CancelBooking) // 取消预订
			}

			// 管理员路由
			admin := authorized.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				// 用户管理
				admin.GET("/users", userHandler.ListUsers)
				admin.GET("/users/:id", userHandler.GetUserByID)
				admin.POST("/users/user", userHandler.AddUser)
				admin.POST("/users/batch", userHandler.DeleteUsers)
				// 预订管理
				admin.GET("/bookings", bookingHandler.ListAllBookings)
				admin.GET("/bookings/search", bookingHandler.SearchBookingsByGuestInfo) // 通过客人信息搜索预订
				admin.POST("/bookings/:id/confirm", bookingHandler.ConfirmBooking)
				admin.POST("/bookings/:id/checkin", bookingHandler.CheckIn)
				admin.POST("/bookings/:id/checkout", bookingHandler.CheckOut)
				admin.GET("/bookings/room", bookingHandler.GetBookingsByRoomNumberAndStatus) // 根据房间号和状态获取预订列表
				// 日志管理
				admin.GET("/logs", logHandler.GetLogs) // 获取日志列表
				// 设施管理
				admin.GET("/facilities", facilityHandler.FindAllFacilities)                  // 查询所有设施
				admin.POST("/facilities", facilityHandler.CreateFacility)                    // 创建设施
				admin.POST("/facilities/batch", facilityHandler.BatchUpdateFacilities)       // 批量更新设施位置
				admin.GET("/facilities/floor/:floor", facilityHandler.FindFacilitiesByFloor) // 按楼层查询设施
				admin.GET("/facilities/:id", facilityHandler.FindFacilityByID)               // 根据ID查找设施
				admin.POST("/facilities/:id", facilityHandler.UpdateFacility)                // 更新设施
				admin.POST("/facilities/:id/delete", facilityHandler.DeleteFacility)         // 删除设施

				// 活动横幅管理
				admin.GET("/banners", bannerHandler.GetAllBanners)            // 获取所有活动横幅
				admin.POST("/banners", bannerHandler.CreateBanner)            // 创建活动横幅
				admin.GET("/banners/:id", bannerHandler.GetBannerByID)        // 获取活动横幅详情
				admin.POST("/banners/:id", bannerHandler.UpdateBanner)        // 更新活动横幅
				admin.POST("/banners/:id/delete", bannerHandler.DeleteBanner) // 删除活动横幅
			}
		}
	}
}
