package middleware

import (
	"net/http"
	"strings"

	"memorylink/internal/config"
	"memorylink/internal/utils"

	"github.com/gin-gonic/gin"
)

// Auth 校验 JWT，将 userId 注入上下文。
func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(utils.ExtractToken(c))
		if token == "" {
			utils.Fail(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(cfg.JWTSecret, token)
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, "登录已失效，请重新登录")
			c.Abort()
			return
		}
		c.Set(utils.CtxUserID, claims.UserID)
		c.Next()
	}
}
