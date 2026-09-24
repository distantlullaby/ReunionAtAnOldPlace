package router

import (
	"memorylink/internal/config"
	"memorylink/internal/handler"
	"memorylink/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup 装配所有路由与依赖。
func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 静态资源：上传的图片。
	r.Static("/uploads", cfg.UploadDir)

	// 组装 handler。
	authH := handler.NewAuthHandler(newAuthSvc(db, cfg))
	storyH := handler.NewStoryHandler(newStorySvc(db))
	userH := handler.NewUserHandler(newUserSvc(db))

	api := r.Group("/api")

	// 无需登录。
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "msg": "pong"})
	})
	api.POST("/auth/register", authH.Register)
	api.POST("/auth/login", authH.Login)
	api.GET("/stories", storyH.List)
	api.GET("/stories/:id", storyH.Detail)

	// 需要登录。
	authed := api.Group("")
	authed.Use(middleware.Auth(cfg))
	{
		authed.GET("/user/profile", userH.Profile)
		authed.PUT("/user/profile", userH.Update)

		authed.POST("/upload", handler.Upload)
		authed.POST("/stories", storyH.Create)
		authed.POST("/stories/:id/append", storyH.AppendReward)
		authed.POST("/stories/:id/respond", storyH.Respond)
		authed.POST("/stories/:id/accept/:rid", storyH.Accept)
	}

	return r
}
