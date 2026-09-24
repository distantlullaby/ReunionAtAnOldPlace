package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"memory-link/config"
	"memory-link/handlers"
	"memory-link/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// 先连接 MySQL 服务器并确保数据库存在
	serverDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
	serverDB, err := gorm.Open(mysql.Open(serverDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接 MySQL (%s:%s): %v", cfg.DBHost, cfg.DBPort, err)
	}
	if err := serverDB.Exec("CREATE DATABASE IF NOT EXISTS `" + cfg.DBName + "` DEFAULT CHARACTER SET utf8mb4").Error; err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&models.User{}, &models.Story{}, &models.Response{}, &models.CoinTransaction{}); err != nil {
		log.Fatalf("建表失败: %v", err)
	}

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatalf("创建上传目录失败: %v", err)
	}

	r := gin.Default()
	r.Use(handlers.CORS())
	r.Static("/uploads", cfg.UploadDir)

	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register(db))
		api.POST("/login", handlers.Login(db))
		api.GET("/stories", handlers.ListStories(db)) // Feed 流公开可看

		auth := api.Group("", handlers.Auth(db))
		{
			auth.GET("/me", handlers.Me())
			auth.POST("/upload", handlers.Upload(cfg.UploadDir))
			auth.POST("/stories", handlers.CreateStory(db))
			auth.POST("/stories/:id/append", handlers.AppendBounty(db))
			auth.POST("/stories/:id/cancel", handlers.CancelStory(db))
			auth.POST("/stories/:id/respond", handlers.Respond(db))
			auth.POST("/responses/:id/accept", handlers.AcceptResponse(db))
			auth.GET("/my/stories", handlers.MyStories(db))
			auth.GET("/my/responses", handlers.MyResponses(db))
			auth.GET("/my/transactions", handlers.MyTransactions(db))
		}
	}

	log.Printf("记忆连接后端已启动: http://localhost:%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
