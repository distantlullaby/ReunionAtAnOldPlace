package service

import (
	"errors"

	"memorylink/internal/config"
	"memorylink/internal/model"
	"memorylink/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// Register 注册账号，事务内创建用户并赠送初始记忆硬币 + 一条赠送流水。
func (s *AuthService) Register(username, password, nickname string) (*model.User, string, error) {
	if len(username) < 3 || len(password) < 6 {
		return nil, "", errors.New("用户名至少3位，密码至少6位")
	}
	if nickname == "" {
		nickname = username
	}

	var existed int64
	s.db.Model(&model.User{}).Where("username = ?", username).Count(&existed)
	if existed > 0 {
		return nil, "", errors.New("用户名已被注册")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := model.User{
		Username:     username,
		PasswordHash: hash,
		Nickname:     nickname,
		CoinBalance:  0,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if _, err := applyCoin(tx, user.ID, s.cfg.InitialCoins, model.TxRegisterGift, nil,
			"注册赠送初始记忆硬币"); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(s.cfg.JWTSecret, user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}
	// 重新读取以拿到赠送后的余额
	_ = s.db.First(&user, user.ID).Error
	return &user, token, nil
}

// Login 校验账号密码并签发 token。
func (s *AuthService) Login(username, password string) (*model.User, string, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}
	if !utils.CheckPassword(user.PasswordHash, password) {
		return nil, "", errors.New("用户名或密码错误")
	}
	token, err := utils.GenerateToken(s.cfg.JWTSecret, user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}
	return &user, token, nil
}
