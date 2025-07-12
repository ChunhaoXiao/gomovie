package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name   string   `form:"name"`
	Movies []*Movie `gorm:"many2many:movie_categories;"`
}
