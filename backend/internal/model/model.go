package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// StringSlice 以 JSON 文本形式存于单个字段中的字符串切片（用于照片 URL 列表）。
type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *StringSlice) Scan(input interface{}) error {
	var data []byte
	switch v := input.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	case nil:
		*s = StringSlice{}
		return nil
	default:
		return errors.New("StringSlice: unsupported scan type")
	}
	if len(data) == 0 {
		*s = StringSlice{}
		return nil
	}
	return json.Unmarshal(data, s)
}

// User 用户表。
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Nickname     string    `gorm:"size:50" json:"nickname"`
	Avatar       string    `gorm:"size:500" json:"avatar"`
	CoinBalance  int       `gorm:"not null;default:0" json:"coinBalance"` // 可用记忆硬币余额
	CreatedAt    time.Time `json:"createdAt"`
}

// 故事状态。
const (
	StoryOpen      = "open"      // 求看中（可被回应/采纳）
	StoryCompleted = "completed" // 已采纳结算
)

// Story 故事表（一条“求看”需求 = 一张双面明信片的“过去”那一面）。
// Reward 是发帖时从发起人余额冻结、托管在故事上的悬赏硬币。
type Story struct {
	ID                 uint        `gorm:"primaryKey" json:"id"`
	UserID             uint        `gorm:"index;not null" json:"userId"`
	Title              string      `gorm:"size:100;not null" json:"title"`
	Location           string      `gorm:"size:100" json:"location"`
	MemoryText         string      `gorm:"type:text" json:"memoryText"`     // 过去回忆
	OldPhotos          StringSlice `gorm:"type:text" json:"oldPhotos"`      // 老照片 URL 列表
	Reward             int         `gorm:"not null;default:0" json:"reward"` // 托管中的悬赏硬币
	Status             string      `gorm:"size:20;index;not null;default:open" json:"status"`
	AcceptedResponseID *uint       `gorm:"index" json:"acceptedResponseId,omitempty"` // 被采纳的回应
	CreatedAt          time.Time   `json:"createdAt"`
	UpdatedAt          time.Time   `json:"updatedAt"`
}

// Response 回应表（一次“替他去拍” = 双面明信片的“当下”那一面）。
type Response struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoryID   uint      `gorm:"index;not null" json:"storyId"`
	UserID    uint      `gorm:"index;not null" json:"userId"`
	NewPhoto  string    `gorm:"size:500" json:"newPhoto"` // 当下现场新照
	Message   string    `gorm:"type:text" json:"message"` // 给发起人的寄语
	CreatedAt time.Time `json:"createdAt"`
}

// 硬币流水类型。
const (
	TxRegisterGift = "register_gift"  // 注册赠送：+
	TxFreeze       = "freeze"         // 发帖冻结悬赏：-
	TxAppendFreeze = "append_freeze"  // 追加悬赏冻结：-
	TxRewardIncome = "reward_income"  // 被采纳获得悬赏：+
	TxRefund       = "refund"         // 解冻退回：+（预留）
)

// CoinTransaction 硬币流水表，一条记录对应某用户一次余额变动。
type CoinTransaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"index;not null" json:"userId"`
	Type            string    `gorm:"size:30;not null" json:"type"`
	Amount          int       `gorm:"not null" json:"amount"` // 带符号：正入账、负出账
	BalanceAfter    int       `gorm:"not null" json:"balanceAfter"`
	RelatedStoryID  *uint     `gorm:"index" json:"relatedStoryId,omitempty"`
	Remark          string    `gorm:"size:200" json:"remark"`
	CreatedAt       time.Time `json:"createdAt"`
}
