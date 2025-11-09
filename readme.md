
## **酒店管理系统**

```
backend/
├── cmd/server/          # 应用入口
├── internal/            # 私有代码
│   ├── config/         # 配置
│   ├── database/       # 数据库
│   ├── handler/        # HTTP 处理器
│   ├── middleware/     # 中间件
│   ├── models/         # 数据模型
│   ├── repository/     # 数据访问
│   └── service/        # 业务逻辑
└── pkg/                # 公共代码
    ├── errors/         # 错误处理
    └── utils/          # 工具函数

```
### **API 基础信息**
- **Base URL**: `http://localhost:8080`
- **认证方式**: JWT Token (Bearer Token)
- **数据格式**: JSON

---

### **1 认证相关 API（公开）**

| 方法 | 路径 | 说明 | 需要登录 |
|------|------|------|---------|
| POST | `/api/auth/register` | 用户注册 | ❌ |
| POST | `/api/auth/login` | 用户登录 | ❌ |

**注册示例：**
```json
POST /api/auth/register
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "123456",
  "phone": "13800138000",
  "real_name": "张三"
}
```

**登录示例：**
```json
POST /api/auth/login
{
  "username": "testuser",
  "password": "123456"
}
```

---

### **2 房间相关 API**

#### **公开访问（不需要登录）**

| 方法 | 路径 | 说明 | 需要登录 |
|------|------|------|---------|
| GET | `/api/rooms` | 获取所有房间列表（分页） | ❌ |
| GET | `/api/rooms/available` | 获取可用房间列表（分页） | ❌ |
| GET | `/api/rooms/:id` | 获取房间详情 | ❌ |
| GET | `/api/rooms/search/type` | 按房型搜索房间 | ❌ |

**查询参数：**
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认10）
- `type`: 房间类型（搜索时使用）

#### **管理员功能（需要登录）**

| 方法 | 路径 | 说明 | 需要登录 |
|------|------|------|---------|
| POST | `/api/rooms` | 创建房间 | ✅ |
| PUT | `/api/rooms/:id` | 更新房间信息 | ✅ |
| DELETE | `/api/rooms/:id` | 删除房间 | ✅ |

---

### **3 用户相关 API（需要登录）**

| 方法 | 路径 | 说明 | 需要登录 |
|------|------|------|---------|
| GET | `/api/users/profile` | 获取个人信息 | ✅ |
| PUT | `/api/users/profile` | 更新个人信息 | ✅ |
| PUT | `/api/users/password` | 修改密码 | ✅ |

---

### **4 预订相关 API（需要登录）**

| 方法 | 路径 | 说明 | 需要登录 |
|------|------|------|---------|
| POST | `/api/bookings` | 创建预订 | ✅ |
| GET | `/api/bookings/my` | 获取我的预订列表 | ✅ |
| GET | `/api/bookings/:id` | 获取预订详情 | ✅ |
| PUT | `/api/bookings/:id/cancel` | 取消预订 | ✅ |

**创建预订示例：**
```json
POST /api/bookings
Headers: Authorization: Bearer <your_token>
{
  "room_id": 1,
  "check_in": "2024-12-01",
  "check_out": "2024-12-05",
  "guest_name": "张三",
  "guest_phone": "13800138000",
  "guest_id_card": "110101199001011234",
  "special_request": "需要加床"
}
```

---

### **5管理员 API（需要登录 + 管理员权限）**

#### **用户管理**

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/users` | 获取用户列表 |
| GET | `/api/admin/users/:id` | 获取指定用户信息 |

#### **预订管理**

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/bookings` | 获取所有预订列表 |
| PUT | `/api/admin/bookings/:id/confirm` | 确认预订 |
| PUT | `/api/admin/bookings/:id/checkin` | 办理入住 |
| PUT | `/api/admin/bookings/:id/checkout` | 办理退房 |

---

### **6其他 API**

| 方法 | 路径 | 说明 | 需要登录 |
|------|------|------|---------|
| GET | `/health` | 健康检查 | ❌ |

---

### **认证说明**

对于需要登录的 API，需要在请求头中添加：
```
Authorization: Bearer <your_jwt_token>
```

获取 token：登录后会返回 token，保存后用于后续请求。

---

### **响应格式**

**成功响应：**
```json
{
  "success": true,
  "data": { ... },
  "message": "操作成功"
}
```

**错误响应：**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "资源不存在"
  }
}
```

**分页响应：**
```json
{
  "success": true,
  "data": [ ... ],
  "page": {
    "page": 1,
    "page_size": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

---

