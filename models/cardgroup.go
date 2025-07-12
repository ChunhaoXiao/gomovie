package models

import "gorm.io/gorm"

type Cardgroup struct {
	gorm.Model
	GroupName string `form:"groupName"`
	Cards     []Card
}
