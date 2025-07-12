package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	CardNumber  string `form:"cardNumber"`
	CoinValue   int    `form:"coinValue"`
	CardgroupID uint
}
