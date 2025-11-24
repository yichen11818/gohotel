package main

import (
	"fmt"
	"log"

	"gohotel/internal/config"
	"gohotel/internal/database"
	"gohotel/internal/handler"
	"gohotel/internal/middleware"
	"gohotel/internal/repository"
	"gohotel/internal/service"
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

	// 2. 连接数据库
	fmt.Println("🔌 正在连接数据库...")
	if err := database.InitMySQL(); err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer database.CloseDB()

	// 3. 自动迁移数据库表
	fmt.Println("🔄 正在执行数据库迁移...")
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 4. 插入测试数据（可选）
	if err := database.SeedData(); err != nil {
		log.Fatal("测试数据插入失败:", err)
	}

	// 5. 初始化雪花算法节点
	fmt.Println("❄️  正在初始化雪花算法节点...")
	// 节点ID可以从配置文件读取，这里暂时使用固定值 1
	// 如果是分布式部署，需要确保每个实例使用不同的节点ID（0-1023）
	if err := utils.InitSnowflake(1); err != nil {
		log.Fatal("雪花算法初始化失败:", err)
	}
	fmt.Println("✅ 雪花算法初始化成功!")

	// 6. 初始化依赖注入
	// Repository 层
	userRepo := repository.NewUserRepository(database.DB)
	roomRepo := repository.NewRoomRepository(database.DB)
	bookingRepo := repository.NewBookingRepository(database.DB)

	// Service 层
	userService := service.NewUserService(userRepo)
	roomService := service.NewRoomService(roomRepo)
	bookingService := service.NewBookingService(bookingRepo, roomRepo, userRepo)

	// Handler 层
	userHandler := handler.NewUserHandler(userService)
	roomHandler := handler.NewRoomHandler(roomService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// 7. 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 8. 创建 Gin 引擎
	r := gin.New()

	// 9. 使用中间件
	r.Use(gin.Recovery())                // 恢复中间件（处理 panic）
	r.Use(middleware.CORSMiddleware())   // 跨域中间件
	r.Use(middleware.LoggerMiddleware()) // 日志中间件

	// 10. 设置路由
	setupRoutes(r, userHandler, roomHandler, bookingHandler)

	// 11. 启动服务器
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("🏨 酒店管理系统 API 服务器")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Printf("📍 服务器地址: http://%s\n", config.AppConfig.Server.Port)
	fmt.Printf("📝 运行模式: %s\n", config.AppConfig.Server.Mode)
	fmt.Printf("📚 Swagger 文档: http://%s/swagger/index.html\n", config.AppConfig.Server.Port)
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("API 文档:")
	fmt.Println("  POST   /api/auth/register      - 用户注册")
	fmt.Println("  POST   /api/auth/login         - 用户登录")
	fmt.Println("  GET    /api/rooms              - 获取房间列表")
	fmt.Println("  GET    /api/rooms/:id          - 获取房间详情")
	fmt.Println("  POST   /api/bookings           - 创建预订（需登录）")
	fmt.Println("  GET    /api/bookings/my        - 我的预订（需登录）")
	fmt.Println("═══════════════════════════════════════════════")

	if err := r.Run(config.AppConfig.Server.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

// setupRoutes 设置所有路由
func setupRoutes(r *gin.Engine, userHandler *handler.UserHandler, roomHandler *handler.RoomHandler, bookingHandler *handler.BookingHandler) {
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
				roomsAuth.POST("", roomHandler.CreateRoom)            // 创建房间
				roomsAuth.POST("/:id", roomHandler.UpdateRoom)        // 更新房间
				roomsAuth.POST("/:id/delete", roomHandler.DeleteRoom) // 删除房间
			}
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
			}
		}
	}
}
