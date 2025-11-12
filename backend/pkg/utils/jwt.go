package utils

import (
	"gohotel/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT 令牌
// 参数：用户ID、用户名
// 返回：令牌字符串、错误
func GenerateToken(userID uint, username string) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(config.AppConfig.JWT.ExpireTime)

	// 创建声明
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gohotel",
		},
	}

	// 使用 HS256 算法创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析 JWT 令牌
// 参数：令牌字符串
// 返回：声明、错误
func ParseToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

















