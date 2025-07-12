package models

import "gorm.io/gorm"

type Picture []string

type Actor struct {
	gorm.Model
	Name        string  `form:"name"`
	IsRecommand uint    `form:"isRecommand"`
	Pictures    Picture `gorm:"serializer:json"`
	Movies      []Movie
}
