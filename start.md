
## 📚 **第一部分：Go 语言基础概念**

### 1. **什么是 Go？**
Go 是一门**编译型**、**静态类型**的编程语言，由 Google 开发。

- **编译型**：代码需要先编译成可执行文件再运行（不像 Python 那样直接运行）
- **静态类型**：变量的类型在编译时就确定了

### 2. **Go 中的"面向对象"**
**重要**：Go 没有传统的"类"（class），但可以实现面向对象的思想！

#### 传统面向对象（如 Java）：
```java
class User {
    private String name;
    public void sayHello() { }
}
```

#### Go 的方式：
```go
type User struct {    // struct 是结构体，相当于"类"
    Name string       // 字段，相当于"属性"
}

func (u *User) SayHello() {  // 方法，绑定在 User 上
    fmt.Println("Hello", u.Name)
}
```

**核心区别**：
- Go 用 `struct`（结构体）代替 `class`
- Go 用 `func (接收者) 方法名()` 来定义方法
- Go 没有继承，但有**组合**和**接口**

---

## 🏗️ **第二部分：项目整体架构**

这个项目是一个**酒店管理系统**，采用经典的**三层架构**：

```
用户请求 
   ↓
Handler（处理器层）    ← 接收 HTTP 请求，返回 HTTP 响应
   ↓
Service（业务逻辑层）   ← 处理业务规则（如：房间是否可预订）
   ↓
Repository（数据访问层） ← 与数据库交互（增删改查）
   ↓
数据库（MySQL）
```

### **为什么要分层？**
- **职责清晰**：每层只做一件事
- **易于维护**：修改数据库不影响业务逻辑
- **可测试**：每层可以独立测试

---

## 📁 **第三部分：项目目录结构详解**

让我详细解释每个目录的作用：

```
gohotel/
├── backend/               # 后端代码
│   ├── cmd/              # 命令/入口目录
│   │   └── server/       
│   │       └── main.go   # 🚀 程序入口（启动文件）
│   │
│   ├── internal/         # 私有代码（只能本项目用）
│   │   ├── config/       # 📝 配置管理（数据库地址、端口等）
│   │   ├── database/     # 💾 数据库连接和初始化
│   │   ├── models/       # 📦 数据模型（User、Room、Booking）
│   │   ├── repository/   # 🗄️ 数据访问层（与数据库交互）
│   │   ├── service/      # 🧠 业务逻辑层（处理业务规则）
│   │   ├── handler/      # 🎯 HTTP 处理器（接收请求）
│   │   └── middleware/   # 🔐 中间件（认证、日志等）
│   │
│   ├── pkg/              # 公共代码（可被外部使用）
│   │   ├── errors/       # ❌ 错误处理
│   │   └── utils/        # 🛠️ 工具函数
│   │
│   ├── go.mod            # Go 依赖管理文件
│   └── go.sum            # 依赖版本锁定文件
│
└── frontend/             # 前端代码（如果有）
```

### **重要概念：internal vs pkg**
- `internal/`：私有代码，只能本项目（gohotel）内部使用
- `pkg/`：公共代码，可以被其他项目导入使用

---

## 🔍 **第四部分：代码详细讲解**

### **1. Models（模型层）**

模型定义了数据的结构，看看用户模型：

```9:21:backend/internal/models/user.go
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`                    // 主键
	Username  string    `gorm:"unique;not null;size:50" json:"username"` // 用户名（唯一）
	Email     string    `gorm:"unique;not null;size:100" json:"email"`   // 邮箱（唯一）
	Password  string    `gorm:"not null;size:255" json:"-"`              // 密码（不返回给前端）
	Phone     string    `gorm:"size:20" json:"phone"`                    // 手机号
	RealName  string    `gorm:"size:50" json:"real_name"`                // 真实姓名
	Avatar    string    `gorm:"size:255" json:"avatar"`                  // 头像 URL
	Role      string    `gorm:"default:'user';size:20" json:"role"`      // 角色：user, admin
	Status    string    `gorm:"default:'active';size:20" json:"status"`  // 状态：active, blocked
	CreatedAt time.Time `json:"created_at"`                              // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                              // 更新时间
}
```

**详细解释**：

1. **`type User struct`**：定义一个名为 `User` 的结构体（相当于类）

2. **字段定义**：`字段名 类型 标签`
   - `ID uint`：ID 字段，类型是无符号整数
   - `Username string`：用户名字段，类型是字符串

3. **标签（Tag）**：反引号中的部分
   ```go
   `gorm:"primaryKey" json:"id"`
   ```
   - `gorm:"primaryKey"`：告诉 GORM（数据库ORM）这是主键
   - `json:"id"`：转成 JSON 时字段名是 "id"
   - `json:"-"`：转 JSON 时忽略这个字段（密码不返回）

4. **方法（类似成员函数）**：

```30:33:backend/internal/models/user.go
// IsAdmin 判断是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}
```

- `func (u *User)`：这是 User 的方法，`u` 是接收者（相当于 this）
- 可以这样调用：`user.IsAdmin()`

### **2. Repository（数据访问层）**

Repository 负责与数据库交互：

```9:17:backend/internal/repository/user_repository.go
// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}
```

**详细解释**：

1. **`UserRepository` 结构体**：
   - 有一个字段 `db`，类型是 `*gorm.DB`（数据库连接）
   - `*` 表示指针类型

2. **构造函数**：
   ```go
   func NewUserRepository(db *gorm.DB) *UserRepository {
       return &UserRepository{db: db}
   }
   ```
   - Go 没有构造函数，习惯用 `New...` 函数创建实例
   - `&` 取地址符，返回指针

3. **方法示例**：

```19:22:backend/internal/repository/user_repository.go
// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
```

- `(r *UserRepository)`：接收者是指针，可以修改结构体
- `error` 是返回类型，表示可能出错
- `r.db.Create(user)`：调用 GORM 创建用户

### **3. Service（业务逻辑层）**

Service 处理业务规则，让我看一个文件：

[1 tool called]

好的，看到 Service 层了！让我继续讲解：

**Service 层的关键点**：

```12:20:backend/internal/service/user_service.go
// UserService 用户业务逻辑层
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}
```

**重要概念：依赖注入**
- `UserService` 依赖 `UserRepository`
- 通过构造函数传入（不是在内部创建）
- 这样便于测试和维护

看一个完整的业务逻辑（注册功能）：

```43:86:backend/internal/service/user_service.go
// Register 用户注册
func (s *UserService) Register(req *RegisterRequest) (*models.User, error) {
	// 1. 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, errors.NewDatabaseError("check username", err)
	}
	if exists {
		return nil, errors.NewConflictError("用户名已存在")
	}

	// 2. 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, errors.NewDatabaseError("check email", err)
	}
	if exists {
		return nil, errors.NewConflictError("邮箱已被使用")
	}

	// 3. 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalServerError("密码加密失败")
	}

	// 4. 创建用户对象
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
		RealName: req.RealName,
		Role:     "user",
		Status:   "active",
	}

	// 5. 保存到数据库
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.NewDatabaseError("create user", err)
	}

	return user, nil
}
```

**流程分析**：
1. **检查业务规则**：用户名、邮箱不能重复
2. **数据处理**：密码加密（安全！）
3. **调用 Repository**：保存到数据库
4. **错误处理**：每步都检查 `err`

**Go 的错误处理**：
```go
result, err := someFunction()
if err != nil {
    // 处理错误
    return nil, err
}
// 使用 result
```
- Go 用返回值返回错误（不是异常）
- 必须显式检查错误

### **4. Handler（处理器层）**

让我看看 Handler：

[1 tool called]

完美！现在看 Handler 层：

```12:20:backend/internal/handler/user_handler.go
// UserHandler 用户控制器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户控制器实例
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
```

Handler 依赖 Service，继续依赖注入的模式。

看一个 HTTP 处理函数：

```22:44:backend/internal/handler/user_handler.go
// Register 用户注册
// POST /api/auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest

	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 2. 调用 Service 层
	user, err := h.userService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 3. 返回成功响应
	utils.SuccessWithMessage(c, "注册成功", gin.H{
		"user": user,
	})
}
```

**详细解释**：

1. **`c *gin.Context`**：
   - Gin 框架的上下文对象
   - 包含请求和响应的所有信息

2. **解析请求**：
   ```go
   var req service.RegisterRequest
   c.ShouldBindJSON(&req)
   ```
   - 把 HTTP Body 中的 JSON 自动转成结构体
   - `&req` 是取地址，传指针

3. **调用 Service**：
   ```go
   user, err := h.userService.Register(&req)
   ```
   - Handler 不处理业务逻辑，只负责接收和返回

4. **返回响应**：
   ```go
   utils.SuccessWithMessage(c, "注册成功", gin.H{"user": user})
   ```
   - 统一的响应格式

---

## 🔄 **第五部分：完整的数据流向**

让我用一个**完整的例子**来串联所有层：

### **场景：用户登录**

**1. 用户发送请求**：
```bash
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456"
}
```

**2. 请求到达 Handler**（`user_handler.go`）：

```46:66:backend/internal/handler/user_handler.go
// Login 用户登录
// POST /api/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest

	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 2. 调用 Service 层
	resp, err := h.userService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// 3. 返回成功响应
	utils.SuccessWithMessage(c, "登录成功", resp)
}
```

**3. Service 处理业务逻辑**（`user_service.go`）：

```88:119:backend/internal/service/user_service.go
// Login 用户登录
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 1. 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewUnauthorizedError("用户名或密码错误")
		}
		return nil, errors.NewDatabaseError("find user", err)
	}

	// 2. 检查账号状态
	if !user.IsActive() {
		return nil, errors.NewForbiddenError("账号已被封禁")
	}

	// 3. 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.NewUnauthorizedError("用户名或密码错误")
	}

	// 4. 生成 JWT 令牌
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.NewInternalServerError("生成令牌失败")
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}
```

**4. Repository 查询数据库**（`user_repository.go`）：

```34:42:backend/internal/repository/user_repository.go
// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
```

**5. 返回响应给用户**：
```json
{
  "success": true,
  "message": "登录成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**流程图**：
```
用户浏览器
    ↓ HTTP POST /api/auth/login
Handler（接收请求，解析 JSON）
    ↓ 调用 userService.Login()
Service（验证业务逻辑，检查密码）
    ↓ 调用 userRepo.FindByUsername()
Repository（执行 SQL 查询）
    ↓ SELECT * FROM users WHERE username = ?
数据库（MySQL）
    ↓ 返回用户数据
Repository → Service → Handler → 用户浏览器
```

---

## 🎯 **第六部分：main.go 启动流程**

让我详细讲解程序是如何启动的：

```17:92:backend/cmd/server/main.go
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

	// 5. 初始化依赖注入
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

	// 6. 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 7. 创建 Gin 引擎
	r := gin.New()

	// 8. 使用中间件
	r.Use(gin.Recovery())                // 恢复中间件（处理 panic）
	r.Use(middleware.CORSMiddleware())   // 跨域中间件
	r.Use(middleware.LoggerMiddleware()) // 日志中间件

	// 9. 设置路由
	setupRoutes(r, userHandler, roomHandler, bookingHandler)

	// 10. 启动服务器
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("🏨 酒店管理系统 API 服务器")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Printf("📍 服务器地址: http://%s\n", config.AppConfig.Server.Port)
	fmt.Printf("📝 运行模式: %s\n", config.AppConfig.Server.Mode)
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
```

**启动步骤详解**：

1. **加载配置**：读取数据库连接信息、服务器端口等
2. **连接数据库**：建立 MySQL 连接
3. **数据库迁移**：自动创建表（users、rooms、bookings）
4. **插入测试数据**：方便测试
5. **依赖注入**：从下往上创建对象
   ```go
   // 先创建 Repository（需要数据库连接）
   userRepo := repository.NewUserRepository(database.DB)
   
   // 再创建 Service（需要 Repository）
   userService := service.NewUserService(userRepo)
   
   // 最后创建 Handler（需要 Service）
   userHandler := handler.NewUserHandler(userService)
   ```
6. **创建 Gin 引擎**：HTTP 服务器框架
7. **注册中间件**：在请求到达 Handler 前执行的逻辑
8. **设置路由**：URL 映射到 Handler
9. **启动服务器**：监听端口，等待请求

---

## 🔐 **第七部分：中间件（Middleware）**

中间件就像"过滤器"，在请求到达 Handler 之前执行。

让我看一个认证中间件：

[1 tool called]

完美！看到中间件了：

```11:45:backend/internal/middleware/auth.go
// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, errors.NewUnauthorizedError("请提供认证令牌"))
			c.Abort()
			return
		}

		// 2. 检查格式：Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, errors.NewUnauthorizedError("令牌格式错误，应为: Bearer <token>"))
			c.Abort()
			return
		}

		// 3. 解析令牌
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.ErrorResponse(c, errors.NewUnauthorizedError("令牌无效或已过期"))
			c.Abort()
			return
		}

		// 4. 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// 5. 继续处理请求
		c.Next()
	}
}
```

**中间件工作原理**：

1. **`gin.HandlerFunc`**：返回一个处理函数
2. **`c.Abort()`**：中断请求，不继续执行
3. **`c.Set()`**：在上下文中存储数据（后面的 Handler 可以用）
4. **`c.Next()`**：继续执行下一个处理器

**使用场景**：
```go
// 在 main.go 中
authorized := api.Group("")
authorized.Use(middleware.AuthMiddleware())  // 使用认证中间件
{
    users.GET("/profile", userHandler.GetProfile)  // 需要登录才能访问
}
```

**请求流程**：
```
用户请求 /api/users/profile
    ↓
AuthMiddleware（检查 Token）
    ↓ (Token 有效，设置 user_id)
UserHandler.GetProfile（获取个人信息）
    ↓
返回响应
```

---

## 📦 **第八部分：Go 的包管理**

看看 `go.mod` 文件：

```1:12:backend/go.mod
module gohotel

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.17.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.1
)
```

**详细解释**：

1. **`module gohotel`**：模块名（导入路径前缀）
   - 导入时：`import "gohotel/internal/models"`

2. **`go 1.21`**：Go 版本要求

3. **`require`**：依赖的第三方包
   - `gin`：Web 框架
   - `jwt`：JWT 令牌
   - `gorm`：ORM（数据库操作）
   - `godotenv`：读取 `.env` 配置文件
   - `crypto`：密码加密

**如何管理依赖**：
```bash
# 下载依赖
go mod download

# 添加新依赖（自动）
go get github.com/some/package

# 清理无用依赖
go mod tidy
```

---

## 🎓 **第九部分：重要的 Go 概念总结**

### 1. **指针（Pointer）**
```go
var x int = 10
var p *int = &x    // p 是指向 x 的指针，& 取地址

fmt.Println(*p)    // *p 解引用，输出 10
*p = 20            // 修改指针指向的值
fmt.Println(x)     // x 变成 20
```

**为什么用指针？**
- 避免复制大对象（效率）
- 允许函数修改参数

### 2. **接口（Interface）**
```go
type Animal interface {
    Speak() string
}

type Dog struct{}
func (d Dog) Speak() string { return "Woof!" }

type Cat struct{}
func (c Cat) Speak() string { return "Meow!" }
```

**接口实现是隐式的**：
- 只要实现了接口的所有方法，就自动实现了接口
- 不需要像 Java 那样写 `implements`

### 3. **错误处理**
```go
file, err := os.Open("file.txt")
if err != nil {
    // 处理错误
    return err
}
defer file.Close()  // defer 延迟执行（函数返回时执行）
```

### 4. **Goroutine（并发）**
```go
go someFunction()  // 在新的 goroutine 中异步执行
```

本项目中暂未使用，但这是 Go 的核心特性。

---

## 🛠️ **第十部分：如何运行和测试**

### **1. 安装依赖**
```bash
cd /home/yumi/Desktop/gohotel/backend
go mod download
```

### **2. 配置数据库**
创建 `.env` 文件（参考 `config.example`）：
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=gohotel
SERVER_PORT=:8080
```

### **3. 运行项目**
```bash
go run cmd/server/main.go
```

### **4. 测试 API**
```bash
# 注册用户
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test",
    "email": "test@test.com",
    "password": "123456"
  }'

# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test",
    "password": "123456"
  }'
```

---

## 📊 **第十一部分：项目架构图**

让我用图形展示整体架构：

```
┌─────────────────────────────────────────────────┐
│                  客户端（浏览器/App）                │
└───────────────────┬─────────────────────────────┘
                    │ HTTP Request
                    ↓
┌─────────────────────────────────────────────────┐
│              Middleware（中间件层）                │
│  - CORS（跨域）                                   │
│  - Logger（日志）                                 │
│  - Auth（认证）                                   │
└───────────────────┬─────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│              Handler（处理器层）                   │
│  - 接收 HTTP 请求                                 │
│  - 参数验证                                       │
│  - 返回 HTTP 响应                                 │
│                                                 │
│  UserHandler | RoomHandler | BookingHandler    │
└───────────────────┬─────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│              Service（业务逻辑层）                  │
│  - 处理业务规则                                   │
│  - 数据验证                                       │
│  - 错误处理                                       │
│                                                 │
│  UserService | RoomService | BookingService    │
└───────────────────┬─────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│            Repository（数据访问层）                 │
│  - 数据库操作（CRUD）                              │
│  - SQL 查询                                      │
│                                                 │
│  UserRepo | RoomRepo | BookingRepo             │
└───────────────────┬─────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│              Database（MySQL）                   │
│  - users 表                                      │
│  - rooms 表                                      │
│  - bookings 表                                   │
└─────────────────────────────────────────────────┘
```

---

## 💡 **第十二部分：学习建议**

### **对于 Go 新手：**

1. **先理解基础语法**：
   - 变量声明：`var x int = 10` 或 `x := 10`
   - 函数：`func add(a, b int) int { return a + b }`
   - 结构体和方法

2. **理解"面向对象"的 Go 方式**：
   - 用 `struct` 代替 `class`
   - 用方法绑定到结构体
   - 用组合代替继承

3. **掌握错误处理**：
   ```go
   result, err := someFunc()
   if err != nil {
       // 处理错误
   }
   ```

4. **理解指针**：
   - 什么时候用值，什么时候用指针
   - 函数参数传指针避免复制

### **理解本项目：**

1. **从 Models 开始**：先看数据结构
2. **然后看 Repository**：理解如何操作数据库
3. **接着看 Service**：理解业务逻辑
4. **最后看 Handler**：理解 HTTP 请求处理
5. **main.go**：理解如何把所有东西组合起来

### **推荐学习路径**：

1. 📖 阅读代码时，从一个功能的完整流程开始（如：用户登录）
2. 🔍 跟踪一个请求从 Handler → Service → Repository → 数据库
3. ✏️ 尝试添加新功能（如：房间评论功能）
4. 🐛 调试和修复问题

---

## 🎯 **总结**

这个项目是一个**典型的 Go Web 应用**，使用了：

- **Gin**：Web 框架
- **GORM**：ORM（数据库）
- **JWT**：认证
- **三层架构**：Handler → Service → Repository

**核心思想**：
1. **分层**：职责清晰，易于维护
2. **依赖注入**：松耦合，便于测试
3. **错误处理**：显式检查，不用异常
4. **接口和组合**：Go 的"面向对象"

--- ---
### 1. Go 常用标准库 API

Go 语言拥有一个强大且设计精良的标准库，以下是一些你几乎每天都会用到的包和函数：

*   **`fmt`**: 用于格式化输入和输出。
    *   `fmt.Println()`: 打印并换行。
    *   `fmt.Printf()`: 根据格式说明符进行格式化输出。
    *   `fmt.Sprintf()`: 格式化成字符串并返回，而不是打印出来。
    *   `fmt.Errorf()`: 创建一个格式化的错误信息。

*   **`net/http`**: 用于构建 HTTP 服务和客户端。
    *   `http.ListenAndServe()`: 启动一个 HTTP 服务器。
    *   `http.HandleFunc()`: 注册一个处理特定路径请求的函数。
    *   `http.Get()`, `http.Post()`: 作为客户端发送 HTTP 请求。
    *   `http.ResponseWriter`: 在服务器端用来构建 HTTP 响应。
    *   `*http.Request`: 包含了客户端的所有请求信息。

*   **`encoding/json`**: 用于处理 JSON 数据。
    *   `json.Marshal()`: 将 Go 的结构体（struct）或其他数据类型编码成 JSON 字节流。
    *   `json.Unmarshal()`: 将 JSON 字节流解码到 Go 的结构体中。

*   **`os`**: 提供与操作系统交互的功能。
    *   `os.Getenv()`: 获取环境变量。
    *   `os.ReadFile()`: 读取文件全部内容。
    *   `os.WriteFile()`: 将数据写入文件。
    *   `os.Args`: 获取命令行参数。

*   **`io`**: 提供 I/O 操作的基本接口。
    *   `io.Reader`: 所有可读取对象的抽象接口（如文件、HTTP 请求体）。
    *   `io.Writer`: 所有可写入对象的抽象接口（如文件、HTTP 响应体）。
    *   `io.Copy()`: 将数据从 `Reader` 复制到 `Writer`。

*   **`sync`**: 用于并发编程中的同步。
    *   `sync.Mutex`: 互斥锁，用于保护共享资源。
    *   `sync.WaitGroup`: 用于等待一组 goroutine 完成执行。

### 2. Gin 框架常用 API

Gin 是一个非常流行的高性能 Go Web 框架，它的 API 设计简洁。

*   **路由 (Routing)**:
    *   `gin.Default()`: 创建一个带有默认中间件（Logger 和 Recovery）的 Gin 引擎。
    *   `router.GET()`, `router.POST()`, `router.PUT()`, `router.DELETE()`: 定义不同 HTTP 方法的路由。
    *   `router.Group()`: 为一组路由创建一个共享相同前缀或中间件的路由组。

*   **请求处理 (Request Handling)**: `c` 通常是 `*gin.Context` 的实例。
    *   `c.Param("id")`: 获取 URL 路径参数，例如 `/users/:id` 中的 `id`。
    *   `c.Query("name")`: 获取 URL 查询参数，例如 `/search?name=John` 中的 `name`。
    *   `c.PostForm("message")`: 获取表单数据。
    *   `c.ShouldBindJSON(&yourStruct)`: 将请求体中的 JSON 数据绑定到一个 Go 结构体上，这是处理 API 请求最常用的方法。

*   **响应处理 (Response Handling)**:
    *   `c.JSON(http.StatusOK, gin.H{"message": "success"})`: 返回 JSON 响应。`gin.H` 是 `map[string]interface{}` 的快捷方式。
    *   `c.String(http.StatusOK, "hello world")`: 返回纯文本响应。
    *   `c.HTML()`: 渲染 HTML 模板并返回。
    *   `c.Redirect()`: 重定向到另一个 URL。

*   **中间件 (Middleware)**:
    *   `router.Use(yourMiddleware())`: 将中间件应用到所有路由。中间件是处理请求前或响应后的函数，常用于日志、认证、CORS 等。

### 3. 前端常用 Web API

这些 API 由浏览器提供，是现代前端开发的核心。

*   **DOM 操作 (DOM Manipulation)**: 用于与页面元素交互。
    *   `document.getElementById()`: 通过 ID 获取元素。
    *   `document.querySelector()`: 使用 CSS 选择器获取第一个匹配的元素。
    *   `document.querySelectorAll()`: 获取所有匹配的元素。
    *   `element.addEventListener()`: 为元素添加事件监听器（如 `click`, `submit`）。

*   **网络请求 (Networking)**:
    *   **Fetch API**: 现代的网络请求标准。
        *   `fetch('/api/data')`: 发送一个 GET 请求。
        *   `.then(response => response.json())`: 将响应体解析为 JSON。
        *   `fetch('/api/submit', { method: 'POST', body: JSON.stringify(data) })`: 发送 POST 请求。
    *   **XMLHttpRequest (XHR)**: 较老的 API，但在一些旧代码或特定场景中仍会见到。

*   **存储 (Storage)**: 在用户浏览器中存储数据。
    *   `localStorage`: 持久化存储，关闭浏览器后数据依然存在。
        *   `localStorage.setItem('key', 'value')`
        *   `localStorage.getItem('key')`
    *   `sessionStorage`: 会话级别存储，关闭浏览器标签页后数据被清除。

*   **定时器 (Timers)**:
    *   `setTimeout(callback, delay)`: 在指定的延迟后执行一次函数。
    *   `setInterval(callback, interval)`: 每隔指定的时间重复执行函数。

*   **URL 和历史记录 (URL and History)**:
    *   `window.location`: 获取或设置当前页面的 URL。
    *   `history.pushState()`: 在不刷新页面的情况下改变 URL，是单页应用（SPA）路由的基础。

这些是你在开发过程中最常接触到的 API。建议你可以在实际编码中多多实践，加深理解。如果你对其中任何一部分想了解更多细节，随时可以问我！