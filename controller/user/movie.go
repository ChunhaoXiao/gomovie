package user

import (
	"movie/db"
	"movie/models"
	"movie/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StreamVideo(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie
	err := db.DB.Where("id=?", id).First(&movie).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "数据不存在",
		})
		return
	}
	fileName := movie.Filename
	utils.StreamVideo(fileName, c)
}

func Show(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie
	error := db.DB.Where("id=?", id).First(&movie).Error
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "数据不存在",
		})
	}
	c.HTML(http.StatusOK, "user/movie/detail.html", gin.H{
		"movie": movie,
	})
}

func MovieLists(c *gin.Context) {
	categoryId, exists := c.GetQuery("category")
	var movies []models.Movie
	if exists {
		var category models.Category
		cateId, _ := strconv.Atoi(categoryId)
		err := db.DB.Where("id=?", cateId).First(&category).Error
		if err == nil {
			db.DB.Model(&category).Association("Movies").Find(&movies)
		}

	} else {
		db.DB.Find(&movies)
	}

	c.HTML(http.StatusOK, "user/movie/index.html", gin.H{
		"movies": movies,
	})

}
