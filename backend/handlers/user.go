package handlers

import (
	"net/http"

	"memory-link/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /api/my/stories  我发布的求看需求（含回应）
func MyStories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		me := currentUser(c)
		var stories []models.Story
		db.Preload("Responses", func(tx *gorm.DB) *gorm.DB { return tx.Order("created_at ASC") }).
			Preload("Responses.User").
			Where("user_id = ?", me.ID).
			Order("created_at DESC").
			Find(&stories)
		c.JSON(http.StatusOK, gin.H{"stories": stories})
	}
}

// GET /api/my/responses  我替别人拍的回应
func MyResponses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		me := currentUser(c)
		var resps []models.Response
		db.Preload("User").Where("user_id = ?", me.ID).Order("created_at DESC").Find(&resps)
		// 附上对应的需求信息
		type item struct {
			models.Response
			Story models.Story `json:"story"`
		}
		items := make([]item, 0, len(resps))
		for _, r := range resps {
			var s models.Story
			db.Preload("User").First(&s, r.StoryID)
			items = append(items, item{Response: r, Story: s})
		}
		c.JSON(http.StatusOK, gin.H{"responses": items})
	}
}

// GET /api/my/transactions  我的硬币流水
func MyTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		me := currentUser(c)
		var txs []models.CoinTransaction
		db.Where("user_id = ?", me.ID).Order("created_at DESC").Limit(100).Find(&txs)
		c.JSON(http.StatusOK, gin.H{"transactions": txs})
	}
}
