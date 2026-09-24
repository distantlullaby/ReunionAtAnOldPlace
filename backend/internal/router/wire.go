package router

import (
	"memorylink/internal/config"
	"memorylink/internal/service"

	"gorm.io/gorm"
)

// 依赖装配的薄封装，集中在一处构造 service。
func newAuthSvc(db *gorm.DB, cfg *config.Config) *service.AuthService {
	return service.NewAuthService(db, cfg)
}
func newStorySvc(db *gorm.DB) *service.StoryService {
	return service.NewStoryService(db)
}
func newUserSvc(db *gorm.DB) *service.UserService {
	return service.NewUserService(db)
}
