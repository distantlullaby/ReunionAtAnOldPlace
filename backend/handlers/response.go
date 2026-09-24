package handlers

import (
	"errors"
	"net/http"

	"memory-link/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type respondReq struct {
	NewPhoto string `json:"new_photo" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

// POST /api/stories/:id/respond  替他去拍：上传现场新照与寄语
func Respond(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req respondReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请上传现场照片并写下寄语"})
			return
		}
		me := currentUser(c)
		storyID := c.Param("id")
		var story models.Story
		if err := db.First(&story, storyID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "需求不存在"})
			return
		}
		if story.Status != "open" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "该需求已结束"})
			return
		}
		if story.UserID == me.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不能回应自己的需求"})
			return
		}
		var count int64
		db.Model(&models.Response{}).Where("story_id = ? AND user_id = ? AND status = ?", story.ID, me.ID, "pending").Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "你已回应过，请等待发起人确认"})
			return
		}
		resp := models.Response{
			StoryID:  story.ID,
			UserID:   me.ID,
			NewPhoto: req.NewPhoto,
			Message:  req.Message,
			Status:   "pending",
		}
		if err := db.Create(&resp).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "提交失败"})
			return
		}
		db.Preload("User").First(&resp, resp.ID)
		c.JSON(http.StatusOK, gin.H{"response": resp})
	}
}

// POST /api/responses/:id/accept  发起人确认采纳：事务内将冻结悬赏结算给拍摄者
func AcceptResponse(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		me := currentUser(c)
		respID := c.Param("id")
		err := db.Transaction(func(tx *gorm.DB) error {
			var resp models.Response
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&resp, respID).Error; err != nil {
				return errors.New("回应不存在")
			}
			if resp.Status != "pending" {
				return errors.New("该回应已处理过")
			}
			var story models.Story
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&story, resp.StoryID).Error; err != nil {
				return errors.New("需求不存在")
			}
			if story.UserID != me.ID {
				return errors.New("只有发起人才能确认结算")
			}
			if story.Status != "open" {
				return errors.New("该需求已结束")
			}
			// 发起人：解冻并支出悬赏
			if err := tx.Model(&models.User{}).Where("id = ?", me.ID).
				Update("frozen", gorm.Expr("frozen - ?", story.Bounty)).Error; err != nil {
				return err
			}
			// 拍摄者：获得悬赏硬币
			if err := tx.Model(&models.User{}).Where("id = ?", resp.UserID).
				Update("balance", gorm.Expr("balance + ?", story.Bounty)).Error; err != nil {
				return err
			}
			// 更新需求与回应状态
			if err := tx.Model(&story).Update("status", "settled").Error; err != nil {
				return err
			}
			if err := tx.Model(&resp).Update("status", "accepted").Error; err != nil {
				return err
			}
			if err := tx.Model(&models.Response{}).
				Where("story_id = ? AND status = ? AND id <> ?", story.ID, "pending", resp.ID).
				Update("status", "rejected").Error; err != nil {
				return err
			}
			// 双方各记一条流水
			if err := tx.Create(&models.CoinTransaction{
				UserID:  me.ID,
				Amount:  0,
				Type:    "settle_out",
				StoryID: story.ID,
				Remark:  "采纳回应，结算悬赏硬币给拍摄者",
			}).Error; err != nil {
				return err
			}
			return tx.Create(&models.CoinTransaction{
				UserID:  resp.UserID,
				Amount:  story.Bounty,
				Type:    "settle_in",
				StoryID: story.ID,
				Remark:  "回应被采纳，获得悬赏硬币",
			}).Error
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "已确认，悬赏硬币已结算给对方"})
	}
}
