package utils

import (
	"fmt"
	"math/rand"
	"movie/db"
	"movie/models"
)

func SeedMovie() {
	var movieCnt int64
	db.DB.Model(&models.Movie{}).Count(&movieCnt)
	categories := seedCategry()
	actors := seedActor()
	if movieCnt < 1000 {
		seedCnt := 1000 - movieCnt
		var i int64
		for i = 0; i < seedCnt; i++ {
			movieCates := []*models.Category{}
			n1 := rand.Intn(5)
			if n1 == 0 {
				movieCates = append(movieCates, &categories[0])
			} else {
				movieCates = append(movieCates, &categories[n1])
				movieCates = append(movieCates, &categories[n1-1])
			}
			n2 := rand.Intn(7)
			actor := actors[n2]

			str := fmt.Sprintf("%d", i)
			title := "电影电影" + str
			var duration int16 = 1000

			movie := models.Movie{
				Title:      title,
				Filename:   "file path.mp4",
				Thumbnail:  []string{},
				Price:      2,
				Duration:   &duration,
				Categories: movieCates,
				Actor:      actor,
			}
			db.DB.Create(&movie)
		}
	}
}

func seedCategry() []models.Category {
	cateList := []string{
		"分类1", "分类2", "分类3", "分类4", "分类5",
	}
	for _, catename := range cateList {
		var category models.Category
		err := db.DB.Where("name=?", catename).First(&category).Error
		if err != nil {
			newCategory := models.Category{
				Name: catename,
			}
			db.DB.Create(&newCategory)
		}
	}
	var allCategories []models.Category
	db.DB.Where("name in ?", cateList).Find(&allCategories)
	return allCategories
}

func seedActor() []models.Actor {
	actors := []string{
		"演员1", "演员2", "演员3", "演员4", "演员5", "演员5", "演员6", "演员7",
	}

	for _, actorName := range actors {
		var actor models.Actor
		err := db.DB.Where("name=?", actorName).First(&actor).Error
		if err != nil {
			newActor := models.Actor{
				Name:        actorName,
				IsRecommand: 0,
				Pictures:    []string{"aaa", "bbb"},
			}
			db.DB.Create(&newActor)
		}

	}
	var allActors []models.Actor
	db.DB.Where("name in ?", actors).Find(&allActors)
	return allActors
}
