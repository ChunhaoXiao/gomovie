package models

import (
	"gorm.io/gorm"
)

type Thumbnails []string

type Movie struct {
	gorm.Model
	Title      string     `form:"title"`
	Filename   string     `form:"filename"`
	Thumbnail  Thumbnails `gorm:"serializer:json"`
	Price      uint
	Actor      Actor       `gorm:"foreignKey:ActorID"`
	Duration   *int16      `form:"duration"`
	Users      []*User     `gorm:"many2many:user_movies;"`
	Categories []*Category `gorm:"many2many:movie_categories;"`
	ActorID    uint
}

// func OrderStatus(cateId int) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if cateId > 0 {
// 			var category Category
// 			db.Where("id=?", cateId).First(&category)
// 			db.Model(&category).Association("Movies")

// 		}
// 		//return db.Scopes(AmountGreaterThan1000).Where("status IN (?)", status)
// 	}
// }
