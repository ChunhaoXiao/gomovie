package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	CardNumber  string `form:"cardNumber"`
	CoinValue   int    `form:"coinValue"`
	CardgroupID uint
	User        User
	UserID      int
}

func UnusedCard(card string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id is NULL and card_number=?", card)
	}
}
