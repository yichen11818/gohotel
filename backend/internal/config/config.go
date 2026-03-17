package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 应用配置结构体
// 这个结构体包含了应用运行所需的所有配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	COS      COSConfig
	WeChat   WeChatConfig
	Log      LogConfig
}

// WeChatConfig 微信小程序配置
type WeChatConfig struct {
	AppID     string // 小程序 AppID
	AppSecret string // 小程序 AppSecret
}

// COSConfig 腾讯云对象存储配置
type COSConfig struct {
	SecretID   string // 密钥ID
	SecretKey  string // 密钥
	BaseURL    string // 访问域名
	BucketName string // 存储桶名称
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string        // 服务器端口，如 ":8080"
	Mode         string        // 运行模式：debug, release, test
	ReadTimeout  time.Duration // 读取超时时间
	WriteTimeout time.Duration // 写入超时时间
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        // 数据库主机地址
	Port            string        // 数据库端口
	User            string        // 数据库用户名
	Password        string        // 数据库密码
	DBName          string        // 数据库名称
	MaxIdleConns    int           // 最大空闲连接数
	MaxOpenConns    int           // 最大打开连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string        // JWT 签名密钥
	ExpireTime time.Duration // Token 过期时间
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string // 日志级别: debug, info, warn, error
	Filename   string // 日志文件路径
	MaxSize    int    // 单个文件最大大小(MB)
	MaxBackups int    // 保留旧文件数量
	MaxAge     int    // 保留天数
	Compress   bool   // 是否压缩
	Console    bool   // 是否输出到控制台
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string // Redis 主机地址
	Port     string // Redis 端口
	Password string // Redis 密码
	DB       int    // Redis 数据库编号
}

// 全局配置对象
var AppConfig *Config

// Load 加载配置
// 这个函数会从环境变量或 .env 文件中读取配置
func Load() error {
	// 尝试加载 .env 文件（如果存在）
	// 在开发环境中，我们可以创建 .env 文件来存储配置
	_ = godotenv.Load()

	AppConfig = &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", ":8080"),
			Mode:         getEnv("SERVER_MODE", "debug"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "3306"),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "hotel"),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", time.Hour),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpireTime: getDurationEnv("JWT_EXPIRE_TIME", 24*time.Hour),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		COS: COSConfig{
			SecretID:   getEnv("SECRETID", ""),
			SecretKey:  getEnv("SECRETKEY", ""),
			BaseURL:    getEnv("COS_BASE_URL", ""),
			BucketName: getEnv("COS_BUCKET_NAME", ""),
		},
		WeChat: WeChatConfig{
			AppID:     getEnv("WECHAT_APP_ID", ""),
			AppSecret: getEnv("WECHAT_APP_SECRET", ""),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Filename:   getEnv("LOG_FILENAME", "logs/app.log"),
			MaxSize:    getIntEnv("LOG_MAX_SIZE", 100),
			MaxBackups: getIntEnv("LOG_MAX_BACKUPS", 3),
			MaxAge:     getIntEnv("LOG_MAX_AGE", 7),
			Compress:   getEnv("LOG_COMPRESS", "true") == "true",
			Console:    getEnv("LOG_CONSOLE", "true") == "true",
		},
	}

	return nil
}

// GetDSN 获取数据库连接字符串
// DSN (Data Source Name) 是数据库连接的标准格式
// 格式：用户名:密码@tcp(主机:端口)/数据库名?参数
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
	)
}

// ========== 辅助函数：从环境变量读取配置 ==========

// getEnv 获取字符串类型的环境变量
// 如果环境变量不存在，返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getIntEnv 获取整数类型的环境变量
func getIntEnv(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getDurationEnv 获取时间类型的环境变量
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
