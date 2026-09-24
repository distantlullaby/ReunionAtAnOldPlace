package service

import (
	"errors"
	"fmt"

	"memorylink/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StoryService struct {
	db *gorm.DB
}

func NewStoryService(db *gorm.DB) *StoryService {
	return &StoryService{db: db}
}

// CreateStoryInput 发帖入参。
type CreateStoryInput struct {
	Title      string
	Location   string
	MemoryText string
	OldPhotos  model.StringSlice
	Reward     int
}

// CreateStory 发布求看需求：事务内校验余额、冻结并扣减悬赏硬币、托管到故事上。
func (s *StoryService) CreateStory(userID uint, in CreateStoryInput) (*model.Story, error) {
	if in.Title == "" {
		return nil, errors.New("标题不能为空")
	}
	if in.Reward < 0 {
		return nil, errors.New("悬赏硬币不能为负")
	}

	story := model.Story{
		UserID:     userID,
		Title:      in.Title,
		Location:   in.Location,
		MemoryText: in.MemoryText,
		OldPhotos:  in.OldPhotos,
		Reward:     in.Reward,
		Status:     model.StoryOpen,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 先锁用户行、扣冻结款；余额不足会在 applyCoin 内回滚。
		remark := fmt.Sprintf("发布求看《%s》冻结悬赏", in.Title)
		if _, err := applyCoin(tx, userID, -in.Reward, model.TxFreeze, nil, remark); err != nil {
			return err
		}
		return tx.Create(&story).Error
	})
	if err != nil {
		return nil, err
	}
	return &story, nil
}

// AppendReward 追加悬赏：继续从发起人余额冻结并累加到故事托管金额上。
func (s *StoryService) AppendReward(userID, storyID uint, add int) (*model.Story, error) {
	if add <= 0 {
		return nil, errors.New("追加硬币需大于0")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var story model.Story
		// 锁定故事行，防止与采纳结算并发。
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&story, storyID).Error; err != nil {
			return errors.New("故事不存在")
		}
		if story.UserID != userID {
			return errors.New("只能给自己的求看追加悬赏")
		}
		if story.Status != model.StoryOpen {
			return errors.New("该求看已结算，无法追加")
		}

		remark := fmt.Sprintf("追加求看《%s》悬赏冻结", story.Title)
		if _, err := applyCoin(tx, userID, -add, model.TxAppendFreeze, &storyID, remark); err != nil {
			return err
		}
		return tx.Model(&story).Update("reward", story.Reward+add).Error
	})
	if err != nil {
		return nil, err
	}

	var story model.Story
	_ = s.db.First(&story, storyID).Error
	return &story, nil
}

// AddResponse 他人“替他去拍”：上传现场新照与寄语。
func (s *StoryService) AddResponse(userID, storyID uint, newPhoto, message string) (*model.Response, error) {
	if newPhoto == "" {
		return nil, errors.New("请上传当下现场照片")
	}
	var story model.Story
	if err := s.db.First(&story, storyID).Error; err != nil {
		return nil, errors.New("故事不存在")
	}
	if story.Status != model.StoryOpen {
		return nil, errors.New("该求看已结算")
	}
	if story.UserID == userID {
		return nil, errors.New("不能替自己去拍")
	}

	resp := model.Response{
		StoryID:  storyID,
		UserID:   userID,
		NewPhoto: newPhoto,
		Message:  message,
	}
	if err := s.db.Create(&resp).Error; err != nil {
		return nil, err
	}
	return &resp, nil
}

// AcceptResponse 发起人确认某条回应：事务内把托管悬赏结算给拍摄者，故事置为完成。
func (s *StoryService) AcceptResponse(ownerID, storyID, responseID uint) (*model.Story, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var story model.Story
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&story, storyID).Error; err != nil {
			return errors.New("故事不存在")
		}
		if story.UserID != ownerID {
			return errors.New("只能确认自己求看下的回应")
		}
		if story.Status != model.StoryOpen {
			return errors.New("该求看已结算，不能重复确认")
		}

		var resp model.Response
		if err := tx.First(&resp, responseID).Error; err != nil {
			return errors.New("回应不存在")
		}
		if resp.StoryID != storyID {
			return errors.New("回应与求看不匹配")
		}
		if resp.UserID == ownerID {
			return errors.New("不能采纳自己的回应")
		}

		reward := story.Reward
		remark := fmt.Sprintf("求看《%s》被采纳，获得悬赏", story.Title)
		if _, err := applyCoin(tx, resp.UserID, reward, model.TxRewardIncome, &storyID, remark); err != nil {
			return err
		}

		// 托管金额随故事结算而释放，余额清零，故事标记完成并记录采纳的回应。
		return tx.Model(&story).Updates(map[string]interface{}{
			"status":                model.StoryCompleted,
			"accepted_response_id":  resp.ID,
			"reward":                0,
		}).Error
	})
	if err != nil {
		return nil, err
	}

	var story model.Story
	_ = s.db.First(&story, storyID).Error
	return &story, nil
}

// StoryDetail 故事及其全部回应（双面明信片）。
type StoryDetail struct {
	model.Story
	Responses []ResponseAuthor `json:"responses"`
}

type ResponseAuthor struct {
	model.Response
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// ListStories Feed 流：分页返回求看列表，附带发起人信息与回应数。
func (s *StoryService) ListStories(page, size int, status string) ([]map[string]interface{}, int64, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}

	q := s.db.Model(&model.Story{})
	if status == model.StoryOpen || status == model.StoryCompleted {
		q = q.Where("status = ?", status)
	}

	var total int64
	q.Count(&total)

	var stories []model.Story
	if err := q.Order("created_at DESC").
		Offset((page - 1) * size).Limit(size).Find(&stories).Error; err != nil {
		return nil, 0, err
	}

	if len(stories) == 0 {
		return []map[string]interface{}{}, total, nil
	}

	userIDs := map[uint]bool{}
	for _, st := range stories {
		userIDs[st.UserID] = true
	}
	var users []model.User
	s.db.Where("id IN ?", mapKeys(userIDs)).Find(&users)
	userMap := map[uint]model.User{}
	for _, u := range users {
		userMap[u.ID] = u
	}

	storyIDs := make([]uint, 0, len(stories))
	for _, st := range stories {
		storyIDs = append(storyIDs, st.ID)
	}
	type cnt struct {
		StoryID uint
		C       int64
	}
	var cnts []cnt
	s.db.Model(&model.Response{}).
		Select("story_id AS story_id, count(*) AS c").
		Where("story_id IN ?", storyIDs).
		Group("story_id").Scan(&cnts)
	cntMap := map[uint]int64{}
	for _, c := range cnts {
		cntMap[c.StoryID] = c.C
	}

	out := make([]map[string]interface{}, 0, len(stories))
	for _, st := range stories {
		u := userMap[st.UserID]
		out = append(out, map[string]interface{}{
			"story":         st,
			"authorName":    u.Nickname,
			"authorAvatar":  u.Avatar,
			"responseCount": cntMap[st.ID],
		})
	}
	return out, total, nil
}

func mapKeys(m map[uint]bool) []uint {
	keys := make([]uint, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GetStory 取单个故事详情（含回应列表与拍摄者信息）。
func (s *StoryService) GetStory(storyID uint) (*StoryDetail, error) {
	var story model.Story
	if err := s.db.First(&story, storyID).Error; err != nil {
		return nil, errors.New("故事不存在")
	}

	var responses []model.Response
	s.db.Where("story_id = ?", storyID).Order("created_at DESC").Find(&responses)

	out := &StoryDetail{Story: story, Responses: []ResponseAuthor{}}
	if len(responses) > 0 {
		uids := map[uint]bool{}
		for _, r := range responses {
			uids[r.UserID] = true
		}
		var users []model.User
		s.db.Where("id IN ?", mapKeys(uids)).Find(&users)
		um := map[uint]model.User{}
		for _, u := range users {
			um[u.ID] = u
		}
		for _, r := range responses {
			u := um[r.UserID]
			out.Responses = append(out.Responses, ResponseAuthor{
				Response: r,
				Nickname: u.Nickname,
				Avatar:   u.Avatar,
			})
		}
	}
	return out, nil
}
