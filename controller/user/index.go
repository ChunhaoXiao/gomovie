package user

import (
	"movie/db"
	"movie/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	var movies []models.Movie
	db.DB.Model(models.Movie{}).Preload("Actor").Order("id desc").Limit(10).Offset(0).Find(&movies)
	//fmt.Println("index movies:::::::", movies)
	c.HTML(http.StatusOK, "user/home/index.html", gin.H{
		"movies": movies,
	})
}
