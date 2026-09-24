package db

import (
	"fmt"
	"log"
	"time"

	"memorylink/internal/config"
	"memorylink/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库句柄。
var DB *gorm.DB

// Init 建库（如不存在）、连接并执行表结构迁移。
func Init(cfg *config.Config) {
	// 1) 不指定库名连接，确保数据库存在。
	root, err := gorm.Open(mysql.Open(cfg.RootDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}
	createSQL := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci",
		cfg.DBName)
	if err := root.Exec(createSQL).Error; err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}
	if sqlDB, err := root.DB(); err == nil {
		_ = sqlDB.Close()
	}

	// 2) 连接目标库。
	gormDB, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("打开数据库 %s 失败: %v", cfg.DBName, err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("获取 sql.DB 失败: %v", err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := gormDB.AutoMigrate(
		&model.User{},
		&model.Story{},
		&model.Response{},
		&model.CoinTransaction{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	DB = gormDB
	log.Printf("数据库 %s 初始化完成", cfg.DBName)
}
