package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"memory-link/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const RegisterBonus int64 = 100 // 注册赠送的初始记忆硬币

func hashPassword(pw string) string {
	sum := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(sum[:])
}

func genToken() string {
	b := make([]byte, 24)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// CORS 中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// 鉴权中间件：Authorization: Bearer <token>
func Auth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		token := strings.TrimPrefix(h, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}
		var user models.User
		if err := db.Where("token = ?", token).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录已失效"})
			return
		}
		c.Set("user", &user)
		c.Next()
	}
}

func currentUser(c *gin.Context) *models.User {
	return c.MustGet("user").(*models.User)
}

type authReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
}

// POST /api/register 注册即赠送初始硬币
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req authReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数不完整"})
			return
		}
		if req.Nickname == "" {
			req.Nickname = req.Username
		}
		user := models.User{
			Username: req.Username,
			Password: hashPassword(req.Password),
			Nickname: req.Nickname,
			Balance:  RegisterBonus,
			Token:    genToken(),
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
			return tx.Create(&models.CoinTransaction{
				UserID: user.ID,
				Amount: RegisterBonus,
				Type:   "register_bonus",
				Remark: "注册赠送初始记忆硬币",
			}).Error
		})
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": user.Token, "user": user})
	}
}

// POST /api/login
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req authReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数不完整"})
			return
		}
		var user models.User
		if err := db.Where("username = ? AND password = ?", req.Username, hashPassword(req.Password)).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		user.Token = genToken()
		db.Model(&user).Update("token", user.Token)
		c.JSON(http.StatusOK, gin.H{"token": user.Token, "user": user})
	}
}

// GET /api/me
func Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user": currentUser(c)})
	}
}
