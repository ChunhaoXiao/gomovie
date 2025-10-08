package models

import "gorm.io/gorm"

type UserCoin struct {
	gorm.Model
	Value  int
	UserID int
	//User   User
}
