package templatefunc

import (
	"fmt"
	"html/template"
	"movie/db"
	"movie/models"
	"slices"

	"github.com/gin-gonic/gin"
)

func FormatDuration(duration *int16) int16 {
	return *duration / 60
}

//cateids := int[]{}

func IsCategorySelected(id int16) bool {
	var categories []models.Category
	db.DB.Find(&categories)
	cateids := []int16{}
	for _, category := range categories {
		cateids = append(cateids, int16(category.ID))
	}
	return slices.Contains(cateids, id)
}

func GetActors() []models.Actor {
	var actors []models.Actor
	db.DB.Find(&actors)
	return actors
}

func GetCategories() []models.Category {
	var categories []models.Category
	db.DB.Find(&categories)
	return categories
}

func ActorName(movie models.Movie) string {
	return movie.Actor.Name
}

func UserIndexActor() []models.Actor {
	var indexActors []models.Actor
	subQuery := db.DB.Model(&models.Movie{}).Select("actor_id", "COUNT(id) as movie_count").Group("actor_id")
	db.DB.Joins("LEFT JOIN (?) as counts on actors.id=counts.actor_id", subQuery).Order("counts.movie_count DESC").Find(&indexActors)
	fmt.Println("index actors::::", indexActors)
	//db.DB.Order("is_recommand desc").Offset(0).Limit(10).Find(&indexActors)
	return indexActors
}

func Islogined() string {
	var c *gin.Context
	cvalue, err := c.Cookie("user")
	if err == nil {
		return cvalue
	}
	return "ff"
}

func GetChargeConfig() template.HTML {
	var data models.Configuration
	err := db.DB.Where("key_name=?", "charge").First(&data).Error
	if err != nil {
		return ""
	}
	//template.HTML(data.Value)
	return template.HTML(data.Value)

}

// func containsPerson(categories []models.Category, target models.Category) bool {
// 	for _, category := range categories {
// 		if category == target {
// 			return true
// 		}
// 	}
// 	return false
// }
