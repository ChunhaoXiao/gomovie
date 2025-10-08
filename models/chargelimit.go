package models

import "gorm.io/gorm"

type ChargeLimit struct {
	gorm.Model
	UserId   int
	TryTimes int
}
