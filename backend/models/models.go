package models

import "time"

// 用户表
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64" json:"username"`
	Password  string    `gorm:"size:128" json:"-"`
	Nickname  string    `gorm:"size:64" json:"nickname"`
	Balance   int64     `gorm:"default:0" json:"balance"` // 可用记忆硬币
	Frozen    int64     `gorm:"default:0" json:"frozen"`  // 冻结中的悬赏硬币
	Token     string    `gorm:"size:64;index" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// 故事表（求看需求）
type Story struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"index" json:"user_id"`
	User       User       `gorm:"foreignKey:UserID" json:"user"`
	Title      string     `gorm:"size:128" json:"title"`
	Location   string     `gorm:"size:128" json:"location"`     // 回忆地点
	MemoryText string     `gorm:"type:text" json:"memory_text"` // 过去回忆
	OldPhoto   string     `gorm:"size:255" json:"old_photo"`    // 老照片
	Bounty     int64      `gorm:"default:0" json:"bounty"`      // 当前冻结的悬赏硬币
	Status     string     `gorm:"size:16;default:open;index" json:"status"` // open / settled / cancelled
	Responses  []Response `gorm:"foreignKey:StoryID" json:"responses,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// 回应表（替他去拍）
type Response struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoryID   uint      `gorm:"index" json:"story_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	NewPhoto  string    `gorm:"size:255" json:"new_photo"` // 现场新照片
	Message   string    `gorm:"type:text" json:"message"`  // 寄语
	Status    string    `gorm:"size:16;default:pending" json:"status"` // pending / accepted / rejected
	CreatedAt time.Time `json:"created_at"`
}

// 硬币流水表
type CoinTransaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Amount    int64     `json:"amount"` // 正=收入，负=支出
	Type      string    `gorm:"size:32" json:"type"` // register_bonus / post_freeze / append_freeze / settle_out / settle_in / cancel_refund
	StoryID   uint      `gorm:"index" json:"story_id"`
	Remark    string    `gorm:"size:255" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}
