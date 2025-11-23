package test

import (
	"encoding/json"
	"gohotel/internal/handler"
	"gohotel/internal/models"
	"gohotel/internal/repository"
	"gohotel/internal/service"
	"gohotel/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 初始化一个内存中的 SQLite 数据库用于测试
func setupTestDB(t *testing.T) *gorm.DB {
	// 使用 :memory: 而不是 file::memory:?cache=shared，确保每个测试都有独立的数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("无法连接到内存数据库: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&models.User{}, &models.Room{}, &models.Booking{})
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	return db
}

// setupRouter 配置用于测试的 Gin 路由
func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 依赖注入
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 设置路由
	api := router.Group("/api/admin")
	{
		api.GET("/users", userHandler.ListUsers)
	}

	return router
}

func TestListUsers_IDSerialization(t *testing.T) {
	// 1. 初始化测试环境
	db := setupTestDB(t)
	router := setupRouter(db)

	// 2. 在数据库中创建一个测试用户，使用一个会丢失精度的大 ID
	largeID := utils.JSONInt64(1989992451309310123)
	testUser := &models.User{
		ID:       largeID,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	result := db.Create(testUser)
	assert.NoError(t, result.Error)

	// 3. 创建一个 HTTP 请求和响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users", nil)

	// 4. 执行请求
	router.ServeHTTP(w, req)

	// 5. 断言响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. 解析响应体并验证 ID 是否为字符串
	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			ID string `json:"id"` // 期望 ID 是一个字符串
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "JSON 反序列化失败")

	// 7. 验证结果
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 1, "应该只返回一个用户")
	assert.Equal(t, largeID.String(), response.Data[0].ID, "响应中的 ID 应该与创建时的 ID 字符串形式匹配")

	// 打印结果以供调试
	t.Logf("响应体: %s", w.Body.String())
	t.Logf("期望的 ID 字符串: %s", largeID.String())
	t.Logf("收到的 ID 字符串: %s", response.Data[0].ID)
}

func TestListUsers_Empty(t *testing.T) {
	// 1. 初始化测试环境
	db := setupTestDB(t)
	router := setupRouter(db)

	// 2. 创建一个 HTTP 请求和响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users", nil)

	// 3. 执行请求
	router.ServeHTTP(w, req)

	// 4. 断言响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. 解析响应体并验证
	var response struct {
		Success bool          `json:"success"`
		Data    []models.User `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "JSON 反序列化失败")

	// 6. 验证结果
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 0, "应该返回一个空的用户列表")

	// 打印结果以供调试
	t.Logf("响应体: %s", w.Body.String())
}

func TestListUsers_WithPagination(t *testing.T) {
	// 1. 初始化测试环境
	db := setupTestDB(t)
	router := setupRouter(db)

	// 2. 在数据库中创建多个测试用户
	for i := 1; i <= 15; i++ {
		user := &models.User{
			ID:       utils.JSONInt64(int64(i)),
			Username: "user" + string(rune('0'+i)),
			Email:    "user" + string(rune('0'+i)) + "@example.com",
			Password: "password",
		}
		result := db.Create(user)
		assert.NoError(t, result.Error)
	}

	// 3. 测试第一页（默认每页10条）
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users?page=1&page_size=10", nil)
	router.ServeHTTP(w, req)

	// 4. 断言响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. 解析响应体并验证
	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			ID string `json:"id"`
		} `json:"data"`
		Page struct {
			Page       int `json:"page"`
			PageSize   int `json:"page_size"`
			Total      int `json:"total"`
			TotalPages int `json:"total_pages"`
		} `json:"page"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "JSON 反序列化失败")

	// 6. 验证结果
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 10, "第一页应该返回10个用户")
	assert.Equal(t, 1, response.Page.Page, "当前页应该是第1页")
	assert.Equal(t, 10, response.Page.PageSize, "每页大小应该是10")
	assert.Equal(t, 15, response.Page.Total, "总用户数应该是15")
	assert.Equal(t, 2, response.Page.TotalPages, "总页数应该是2")

	// 7. 测试第二页
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/admin/users?page=2&page_size=10", nil)
	router.ServeHTTP(w2, req2)

	var response2 struct {
		Success bool `json:"success"`
		Data    []struct {
			ID string `json:"id"`
		} `json:"data"`
		Page struct {
			Page       int `json:"page"`
			PageSize   int `json:"page_size"`
			Total      int `json:"total"`
			TotalPages int `json:"total_pages"`
		} `json:"page"`
	}

	err2 := json.Unmarshal(w2.Body.Bytes(), &response2)
	assert.NoError(t, err2, "JSON 反序列化失败")

	assert.True(t, response2.Success)
	assert.Len(t, response2.Data, 5, "第二页应该返回5个用户")
	assert.Equal(t, 2, response2.Page.Page, "当前页应该是第2页")

	// 打印结果以供调试
	t.Logf("第一页响应体: %s", w.Body.String())
	t.Logf("第二页响应体: %s", w2.Body.String())
}

func TestListUsers_WithDifferentRoles(t *testing.T) {
	// 1. 初始化测试环境
	db := setupTestDB(t)
	router := setupRouter(db)

	// 2. 在数据库中创建不同角色的用户
	users := []models.User{
		{ID: utils.JSONInt64(1), Username: "admin", Email: "admin@example.com", Password: "password", Role: "admin"},
		{ID: utils.JSONInt64(2), Username: "user", Email: "user@example.com", Password: "password", Role: "user"},
		{ID: utils.JSONInt64(3), Username: "guest", Email: "guest@example.com", Password: "password", Role: "guest"},
	}
	for i := range users {
		result := db.Create(&users[i])
		assert.NoError(t, result.Error)
	}

	// 3. 创建一个 HTTP 请求和响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users", nil)

	// 4. 执行请求
	router.ServeHTTP(w, req)

	// 5. 断言响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. 解析响应体并验证
	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Role     string `json:"role"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "JSON 反序列化失败")

	// 7. 验证结果
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 3, "应该返回三个用户")

	// 验证每个用户的角色
	assert.Equal(t, "admin", response.Data[0].Role)
	assert.Equal(t, "user", response.Data[1].Role)
	assert.Equal(t, "guest", response.Data[2].Role)

	// 打印结果以供调试
	t.Logf("响应体: %s", w.Body.String())
}

func TestListUsers_WithDifferentStatuses(t *testing.T) {
	// 1. 初始化测试环境
	db := setupTestDB(t)
	router := setupRouter(db)

	// 2. 在数据库中创建不同状态的用户
	users := []models.User{
		{ID: utils.JSONInt64(1), Username: "active_user", Email: "active@example.com", Password: "password", Status: "active"},
		{ID: utils.JSONInt64(2), Username: "inactive_user", Email: "inactive@example.com", Password: "password", Status: "inactive"},
		{ID: utils.JSONInt64(3), Username: "banned_user", Email: "banned@example.com", Password: "password", Status: "banned"},
	}
	for i := range users {
		result := db.Create(&users[i])
		assert.NoError(t, result.Error)
	}

	// 3. 创建一个 HTTP 请求和响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users", nil)

	// 4. 执行请求
	router.ServeHTTP(w, req)

	// 5. 断言响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. 解析响应体并验证
	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Status   string `json:"status"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "JSON 反序列化失败")

	// 7. 验证结果
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 3, "应该返回三个用户")

	// 验证每个用户的状态
	assert.Equal(t, "active", response.Data[0].Status)
	assert.Equal(t, "inactive", response.Data[1].Status)
	assert.Equal(t, "banned", response.Data[2].Status)

	// 打印结果以供调试
	t.Logf("响应体: %s", w.Body.String())
}

func TestListUsers_VerifyAllFields(t *testing.T) {
	// 1. 初始化测试环境
	db := setupTestDB(t)
	router := setupRouter(db)

	// 2. 在数据库中创建一个包含所有字段的测试用户
	phone := "13800138000"
	testUser := &models.User{
		ID:       utils.JSONInt64(1),
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		Phone:    &phone,
		RealName: "测试用户",
		Avatar:   "https://example.com/avatar.jpg",
		Role:     "admin",
		Status:   "active",
	}
	result := db.Create(testUser)
	assert.NoError(t, result.Error)

	// 3. 创建一个 HTTP 请求和响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users", nil)

	// 4. 执行请求
	router.ServeHTTP(w, req)

	// 5. 断言响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. 解析响应体并验证所有字段
	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			ID       string  `json:"id"`
			Username string  `json:"username"`
			Email    string  `json:"email"`
			Phone    *string `json:"phone"`
			RealName string  `json:"real_name"`
			Avatar   string  `json:"avatar"`
			Role     string  `json:"role"`
			Status   string  `json:"status"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "JSON 反序列化失败")

	// 7. 验证结果
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 1, "应该只返回一个用户")

	// 验证所有字段
	user := response.Data[0]
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotNil(t, user.Phone)
	assert.Equal(t, "13800138000", *user.Phone)
	assert.Equal(t, "测试用户", user.RealName)
	assert.Equal(t, "https://example.com/avatar.jpg", user.Avatar)
	assert.Equal(t, "admin", user.Role)
	assert.Equal(t, "active", user.Status)

	// 打印结果以供调试
	t.Logf("响应体: %s", w.Body.String())
}

func TestListUsers_MultipleLargeIDs(t *testing.T) {
	// 1. 初始化测试环境
	db := setupTestDB(t)
	router := setupRouter(db)

	// 2. 在数据库中创建多个具有大 ID 的测试用户
	largeIDs := []utils.JSONInt64{
		utils.JSONInt64(9007199254740991),    // JavaScript Number.MAX_SAFE_INTEGER
		utils.JSONInt64(9007199254740992),    // 超过 JavaScript 安全整数范围
		utils.JSONInt64(1989992451309310123), // 更大的 ID
	}

	for i, id := range largeIDs {
		user := &models.User{
			ID:       id,
			Username: "user" + string(rune('0'+i)),
			Email:    "user" + string(rune('0'+i)) + "@example.com",
			Password: "password",
		}
		result := db.Create(user)
		assert.NoError(t, result.Error)
	}

	// 3. 创建一个 HTTP 请求和响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users", nil)

	// 4. 执行请求
	router.ServeHTTP(w, req)

	// 5. 断言响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. 解析响应体并验证 ID 是否为字符串
	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "JSON 反序列化失败")

	// 7. 验证结果
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 3, "应该返回三个用户")

	// 验证每个 ID 都是正确的字符串形式
	for i, id := range largeIDs {
		assert.Equal(t, id.String(), response.Data[i].ID, "ID 应该正确序列化为字符串")
	}

	// 打印结果以供调试
	t.Logf("响应体: %s", w.Body.String())
	for i, id := range largeIDs {
		t.Logf("期望的 ID[%d]: %s, 收到的 ID[%d]: %s", i, id.String(), i, response.Data[i].ID)
	}
}

func TestListUsers_MultipleUsers(t *testing.T) {
	// 1. 初始化测试环境
	db := setupTestDB(t)
	router := setupRouter(db)

	// 2. 在数据库中创建多个测试用户，手动指定不同的 ID
	users := []models.User{
		{ID: utils.JSONInt64(1), Username: "user1", Email: "user1@example.com", Password: "password"},
		{ID: utils.JSONInt64(2), Username: "user2", Email: "user2@example.com", Password: "password"},
	}
	for i := range users {
		result := db.Create(&users[i])
		assert.NoError(t, result.Error)
	}

	// 3. 创建一个 HTTP 请求和响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users", nil)

	// 4. 执行请求
	router.ServeHTTP(w, req)

	// 5. 断言响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. 解析响应体并验证
	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "JSON 反序列化失败")

	// 7. 验证结果
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2, "应该返回两个用户")

	// 验证返回的用户数据是否正确
	assert.Equal(t, users[0].Username, response.Data[0].Username)
	assert.Equal(t, users[0].Email, response.Data[0].Email)
	assert.Equal(t, users[1].Username, response.Data[1].Username)
	assert.Equal(t, users[1].Email, response.Data[1].Email)

	// 打印结果以供调试
	t.Logf("响应体: %s", w.Body.String())
}
