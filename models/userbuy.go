package models

import "gorm.io/gorm"

type UserBuy struct {
	gorm.Model
	UserID  uint
	MovieID uint
	User    User
	Movie   Movie
}
