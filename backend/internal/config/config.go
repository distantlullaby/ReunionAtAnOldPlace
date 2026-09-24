package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 全局配置，所有项均可由环境变量覆盖，默认值对接本地开发环境。
type Config struct {
	ServerPort   string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	InitialCoins int // 注册即赠送的初始记忆硬币
	UploadDir    string
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

// Load 读取配置。
func Load() *Config {
	return &Config{
		ServerPort:   getenv("SERVER_PORT", "8080"),
		DBHost:       getenv("DB_HOST", "192.168.5.100"),
		DBPort:       getenv("DB_PORT", "3306"),
		DBUser:       getenv("DB_USER", "root"),
		DBPassword:   getenv("DB_PASSWORD", "123456"),
		DBName:       getenv("DB_NAME", "memory_link"),
		JWTSecret:    getenv("JWT_SECRET", "memory-link-secret-change-me"),
		InitialCoins: getenvInt("INITIAL_COINS", 100),
		UploadDir:    getenv("UPLOAD_DIR", "uploads"),
	}
}

// dsn 拼装 MySQL DSN；withDB=false 时不指定库名（用于建库前连接）。
func (c *Config) dsn(withDB bool) string {
	dbName := ""
	if withDB {
		dbName = c.DBName
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, dbName)
}

// DSN 带库名的连接串。
func (c *Config) DSN() string { return c.dsn(true) }

// RootDSN 不带库名的连接串，用于首次建库。
func (c *Config) RootDSN() string { return c.dsn(false) }
