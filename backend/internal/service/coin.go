package service

import (
	"errors"

	"memorylink/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// applyCoin 在一个已开启的事务内对某用户余额做带符号增减，并写一条流水。
// 通过 SELECT ... FOR UPDATE 锁定用户行，保证并发下余额与流水强一致、不会透支。
// amount 为正表示入账，为负表示出账。
func applyCoin(tx *gorm.DB, userID uint, amount int, txType string, relatedStoryID *uint, remark string) (int, error) {
	var user model.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
		return 0, errors.New("用户不存在")
	}

	newBalance := user.CoinBalance + amount
	if newBalance < 0 {
		return 0, errors.New("记忆硬币余额不足")
	}

	user.CoinBalance = newBalance
	if err := tx.Save(&user).Error; err != nil {
		return 0, err
	}

	ct := model.CoinTransaction{
		UserID:         userID,
		Type:           txType,
		Amount:         amount,
		BalanceAfter:   newBalance,
		RelatedStoryID: relatedStoryID,
		Remark:         remark,
	}
	if err := tx.Create(&ct).Error; err != nil {
		return 0, err
	}
	return newBalance, nil
}
