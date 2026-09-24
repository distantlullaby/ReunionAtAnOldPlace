package handlers

import (
	"errors"
	"net/http"

	"memory-link/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GET /api/stories  Feed 流：所有求看需求（新的在前）
func ListStories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stories []models.Story
		if err := db.Preload("User").
			Preload("Responses", func(tx *gorm.DB) *gorm.DB { return tx.Order("created_at ASC") }).
			Preload("Responses.User").
			Order("created_at DESC").
			Find(&stories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"stories": stories})
	}
}

type createStoryReq struct {
	Title      string `json:"title" binding:"required"`
	Location   string `json:"location" binding:"required"`
	MemoryText string `json:"memory_text" binding:"required"`
	OldPhoto   string `json:"old_photo"`
	Bounty     int64  `json:"bounty" binding:"required,min=1"`
}

// POST /api/stories  发布求看需求：事务内冻结悬赏硬币
func CreateStory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createStoryReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数不完整，悬赏至少 1 枚硬币"})
			return
		}
		me := currentUser(c)
		story := models.Story{
			UserID:     me.ID,
			Title:      req.Title,
			Location:   req.Location,
			MemoryText: req.MemoryText,
			OldPhoto:   req.OldPhoto,
			Bounty:     req.Bounty,
			Status:     "open",
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			// 行锁读取用户，防止并发超扣
			var user models.User
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, me.ID).Error; err != nil {
				return err
			}
			if user.Balance < req.Bounty {
				return errors.New("记忆硬币余额不足")
			}
			if err := tx.Model(&user).Updates(map[string]interface{}{
				"balance": gorm.Expr("balance - ?", req.Bounty),
				"frozen":  gorm.Expr("frozen + ?", req.Bounty),
			}).Error; err != nil {
				return err
			}
			if err := tx.Create(&story).Error; err != nil {
				return err
			}
			return tx.Create(&models.CoinTransaction{
				UserID:  me.ID,
				Amount:  -req.Bounty,
				Type:    "post_freeze",
				StoryID: story.ID,
				Remark:  "发布求看需求，冻结悬赏硬币",
			}).Error
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Preload("User").First(&story, story.ID)
		c.JSON(http.StatusOK, gin.H{"story": story})
	}
}

type appendReq struct {
	Amount int64 `json:"amount" binding:"required,min=1"`
}

// POST /api/stories/:id/append  追加悬赏：事务内追加冻结
func AppendBounty(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req appendReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "追加数量至少为 1"})
			return
		}
		me := currentUser(c)
		storyID := c.Param("id")
		err := db.Transaction(func(tx *gorm.DB) error {
			var story models.Story
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&story, storyID).Error; err != nil {
				return errors.New("需求不存在")
			}
			if story.UserID != me.ID {
				return errors.New("只能给自己的需求追加悬赏")
			}
			if story.Status != "open" {
				return errors.New("该需求已结束，无法追加")
			}
			var user models.User
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, me.ID).Error; err != nil {
				return err
			}
			if user.Balance < req.Amount {
				return errors.New("记忆硬币余额不足")
			}
			if err := tx.Model(&user).Updates(map[string]interface{}{
				"balance": gorm.Expr("balance - ?", req.Amount),
				"frozen":  gorm.Expr("frozen + ?", req.Amount),
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&story).Update("bounty", gorm.Expr("bounty + ?", req.Amount)).Error; err != nil {
				return err
			}
			return tx.Create(&models.CoinTransaction{
				UserID:  me.ID,
				Amount:  -req.Amount,
				Type:    "append_freeze",
				StoryID: story.ID,
				Remark:  "追加悬赏，追加冻结硬币",
			}).Error
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "追加成功"})
	}
}

// POST /api/stories/:id/cancel  取消需求：事务内退回冻结硬币
func CancelStory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		me := currentUser(c)
		storyID := c.Param("id")
		err := db.Transaction(func(tx *gorm.DB) error {
			var story models.Story
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&story, storyID).Error; err != nil {
				return errors.New("需求不存在")
			}
			if story.UserID != me.ID {
				return errors.New("只能取消自己的需求")
			}
			if story.Status != "open" {
				return errors.New("该需求已结束")
			}
			if err := tx.Model(&story).Update("status", "cancelled").Error; err != nil {
				return err
			}
			if err := tx.Model(&models.Response{}).Where("story_id = ? AND status = ?", story.ID, "pending").
				Update("status", "rejected").Error; err != nil {
				return err
			}
			if story.Bounty > 0 {
				if err := tx.Model(&models.User{}).Where("id = ?", me.ID).Updates(map[string]interface{}{
					"balance": gorm.Expr("balance + ?", story.Bounty),
					"frozen":  gorm.Expr("frozen - ?", story.Bounty),
				}).Error; err != nil {
					return err
				}
				return tx.Create(&models.CoinTransaction{
					UserID:  me.ID,
					Amount:  story.Bounty,
					Type:    "cancel_refund",
					StoryID: story.ID,
					Remark:  "取消需求，退回冻结悬赏",
				}).Error
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "已取消，悬赏硬币已退回"})
	}
}
