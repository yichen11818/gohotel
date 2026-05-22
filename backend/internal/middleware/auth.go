package middleware

import (
	stderrors "errors"
	"gohotel/internal/database"
	"gohotel/internal/models"
	"gohotel/pkg/errors"
	"gohotel/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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
			if stderrors.Is(err, jwt.ErrTokenExpired) {
				utils.ErrorResponse(c, errors.NewUnauthorizedError("令牌已过期"))
				c.Abort()
				return
			}
			if stderrors.Is(err, jwt.ErrSignatureInvalid) {
				utils.ErrorResponse(c, errors.NewUnauthorizedError("令牌签名无效"))
				c.Abort()
				return
			}
			utils.ErrorResponse(c, errors.NewUnauthorizedError("令牌无效"))
			c.Abort()
			return
		}

		// 4. 从数据库重新确认用户状态与 token 版本，避免密码修改后旧令牌继续生效。
		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			utils.ErrorResponse(c, errors.NewUnauthorizedError("登录状态已失效"))
			c.Abort()
			return
		}
		if !user.IsActive() {
			utils.ErrorResponse(c, errors.NewForbiddenError("账号已被封禁"))
			c.Abort()
			return
		}
		if claims.TokenVersion != user.TokenVersion {
			utils.ErrorResponse(c, errors.NewUnauthorizedError("登录状态已失效，请重新登录"))
			c.Abort()
			return
		}

		// 5. 将最新的用户信息存入上下文
		c.Set("user_id", user.ID.Int64())
		c.Set("username", user.Username)
		c.Set("role", user.Role)

		// 6. 继续处理请求
		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
// 注意：必须在 AuthMiddleware 之后使用
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleAny, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, errors.NewUnauthorizedError("未登录"))
			c.Abort()
			return
		}

		role, ok := roleAny.(string)
		if !ok || role == "" {
			utils.ErrorResponse(c, errors.NewUnauthorizedError("无效的登录信息"))
			c.Abort()
			return
		}

		if role != "admin" {
			utils.ErrorResponse(c, errors.NewForbiddenError("需要管理员权限"))
			c.Abort()
			return
		}

		c.Next()
	}
}
