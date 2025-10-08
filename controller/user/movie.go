package user

import (
	"fmt"
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
	fmt.Println("movie file name is:", fileName)
	utils.StreamVideo(fileName, c)
}

func Show(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie
	var movies []models.Movie
	error := db.DB.Where("id=?", id).Preload("Categories").First(&movie).Error
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "数据不存在",
		})
		return
	}

	fmt.Println("movie categories#########", movie.Categories)

	user, exists := c.Get("currentUser")
	if exists {
		currentUser := user.(models.User)
		//db.DB.Model(&currentUser).Association("Movies").Find()
		//db.DB.Model(&currentUser).Where("movies.id=?", id).Association("Movies").Find(&movies)
		//fmt.Println("bougnt movies", movies)
		// var users []models.User
		// db.DB.Model(&models.User{}).Where("id=?", currentUser.ID).Preload("Movies", "id", id).Find(&users)
		// fmt.Println("moviesssss", users[0].Movies)
		db.DB.Model(&currentUser).Where("id=?", id).Association("Movies").Find(&movies)
	}

	fmt.Println("##############", len(movies))
	c.HTML(http.StatusOK, "user/movie/detail.html", gin.H{
		"movie":     movie,
		"usermovie": len(movies),
	})
}

func MovieLists(c *gin.Context) {
	categoryId := c.DefaultQuery("category", "")
	var movies []models.Movie
	perPage := 30
	if categoryId != "" {
		var category models.Category
		cateId, _ := strconv.Atoi(categoryId)
		err := db.DB.Where("id=?", cateId).First(&category).Error
		if err == nil {
			db.DB.Model(&category).Scopes(utils.Paginate(c.Request, perPage)).Association("Movies").Find(&movies)
		}

	} else {
		db.DB.Scopes(utils.Paginate(c.Request, perPage)).Order("id desc").Find(&movies)
	}

	var movieCount int64
	db.DB.Model(&models.Movie{}).Count(&movieCount)
	//pages := movieCount / int64(perPage)
	currentPageStr := c.DefaultQuery("page", "1")
	currentPage, _ := strconv.Atoi(currentPageStr)

	/*
		var pageList []int
		start := 1
		end := 5
		if pages <= 5 {
			for i := start; i <= int(pages); i++ {
				pageList = append(pageList, i)
			}
		} else {
			if currentPage > 3 {
				start = currentPage - 2
				end = currentPage + 2
			}
			for i := start; i <= end; i++ {
				pageList = append(pageList, i)
			}
		}

		fmt.Println(pageList)
		var prevPage int
		var nextPage int
		if currentPage > 1 {
			prevPage = currentPage - 1
		} else {
			prevPage = 1
		}

		if currentPage < int(pages) {
			nextPage = currentPage + 1
		} else {
			nextPage = currentPage
		} */

	var prevPage int
	var nextPage int
	if currentPage > 1 {
		prevPage = currentPage - 1
	} else {
		prevPage = 1
	}

	if currentPage < int(movieCount/int64(perPage)) {
		nextPage = currentPage + 1
	} else {
		nextPage = currentPage
	}

	totalPages := movieCount / int64(perPage)
	fmt.Println("total Pages", totalPages)

	category := c.DefaultQuery("category", "")
	pageStr := utils.GeneratePageStr(perPage, currentPage, int(movieCount), 5)

	c.HTML(http.StatusOK, "user/movie/index.html", gin.H{
		"movies":      movies,
		"currentPage": currentPage,
		"pages":       pageStr,
		"nextPage":    nextPage,
		"prevPage":    prevPage,
		"category":    category,
	})

}
