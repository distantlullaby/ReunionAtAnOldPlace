package service

import (
	"memorylink/internal/model"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Profile 个人中心汇总：用户信息 + 各类相关记录。
func (s *UserService) Profile(userID uint) (map[string]interface{}, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	var myStories []model.Story
	s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&myStories)

	var myResponses []model.Response
	s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&myResponses)

	// 我的回应附带对应故事标题，便于展示。
	type respWithStory struct {
		model.Response
		StoryTitle string `json:"storyTitle"`
		Status     string `json:"status"`
	}
	respList := []respWithStory{}
	if len(myResponses) > 0 {
		sids := make([]uint, 0, len(myResponses))
		for _, r := range myResponses {
			sids = append(sids, r.StoryID)
		}
		var stories []model.Story
		s.db.Where("id IN ?", sids).Find(&stories)
		sm := map[uint]model.Story{}
		for _, st := range stories {
			sm[st.ID] = st
		}
		for _, r := range myResponses {
			if st, ok := sm[r.StoryID]; ok {
				respList = append(respList, respWithStory{
					Response:   r,
					StoryTitle: st.Title,
					Status:     st.Status,
				})
			}
		}
	}

	var txs []model.CoinTransaction
	s.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(100).Find(&txs)

	return map[string]interface{}{
		"user":          user,
		"myStories":     myStories,
		"myResponses":   respList,
		"transactions":  txs,
		"stats": map[string]interface{}{
			"storyCount":    len(myStories),
			"responseCount": len(myResponses),
		},
	}, nil
}

// UpdateProfile 更新昵称/头像。
func (s *UserService) UpdateProfile(userID uint, nickname, avatar string) error {
	updates := map[string]interface{}{}
	if nickname != "" {
		updates["nickname"] = nickname
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}
	if len(updates) == 0 {
		return nil
	}
	return s.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}
