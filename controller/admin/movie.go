package admin

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"movie/db"
	"movie/dto"
	"movie/models"

	"movie/utils"

	"github.com/gin-gonic/gin"
)

func MovieList(c *gin.Context) {
	var movies []models.Movie
	db.DB.Preload("Actor").Find(&movies)
	//fmt.Println("movie list:", movies)
	for _, obj := range movies {
		fmt.Println(obj.Actor.Name)
	}
	c.HTML(http.StatusOK, "admin/movie/index.html", gin.H{
		"movies": movies,
	})
}

func Create(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/movie/create.tmpl", gin.H{})
}

func SaveMovie(c *gin.Context) {
	var movieCreateDto dto.MovieCreateDto
	var categories []*models.Category
	err := c.ShouldBind(&movieCreateDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "validate failed",
		})
		return
	}
	fmt.Println(movieCreateDto)
	categoriesids := movieCreateDto.Categories
	fmt.Println("categoryIDS:", categoriesids)
	db.DB.Where("id IN ?", categoriesids).Find(&categories)
	thumb := utils.MakeMovieThumb(movieCreateDto.Filename)

	fmt.Println("###cates:::", categories)
	entity := models.Movie{
		Title: movieCreateDto.Title,
		//Actor:      &movieCreateDto.Actor,
		Filename:   movieCreateDto.Filename,
		Duration:   &movieCreateDto.Duration,
		Categories: categories,
		ActorID:    uint(movieCreateDto.ActorId),
		Thumbnail:  &thumb,
	}
	erros := db.DB.Create(&entity).Error
	if erros != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": erros.Error(),
		})
	}
	//db.DB.Model(&entity).Association("Categories").Append(categories)
	c.Redirect(http.StatusMovedPermanently, "/admin/movie/index")

}

func PlayMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "参数错误",
		})
	}
	var movie models.Movie
	errs := db.DB.Where("id=?", id).First(&movie).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "数据不存在",
		})
	}
	utils.StreamVideo(movie.Filename, c)

}

func RemoveMovie(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie
	db.DB.Where("id=?", id).First(&movie)
	os.Remove("./uploads/" + movie.Filename)
	db.DB.Model(&movie).Association("Categories").Clear()
	db.DB.Delete(&movie)

	c.JSON(http.StatusOK, gin.H{
		"data": "success",
	})

}

func EditMovie(c *gin.Context) {
	id := c.Param("id")
	var categories []models.Category
	db.DB.Find(&categories)
	var movie models.Movie
	db.DB.Model(models.Movie{}).Preload("Categories").Where("id=?", id).First(&movie)
	fmt.Println("movie is:", movie.Thumbnail)
	movieCategories := make(map[int]string)
	for _, category := range movie.Categories {
		movieCategories[int(category.ID)] = category.Name
	}
	fmt.Println("movieCategories===>", movieCategories)

	c.HTML(http.StatusOK, "admin/movie/create.tmpl", gin.H{
		"categories": categories,
		"movieId":    movie.ID,
		"title":      movie.Title,
		"duration":   movie.Duration,
		//"actor":           movie.Actor,
		"MovieCategories": movieCategories,
		"filename":        movie.Filename,
		"thumbnail":       movie.Thumbnail,
	})
}

func UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie
	err := db.DB.First(&movie, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "数据不存在",
		})
	}
	var movieCreateDto dto.MovieCreateDto
	errs := c.ShouldBind(&movieCreateDto)
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "validate err",
		})
		return
	}

	var categories []*models.Category
	categoriesids := movieCreateDto.Categories
	if len(categoriesids) > 0 {
		db.DB.Where("id IN ?", categoriesids).Find(&categories)
	}
	thumb := utils.MakeMovieThumb(movieCreateDto.Filename)
	movie.Title = movieCreateDto.Title
	//movie.Actor = &movieCreateDto.Actor
	movie.Filename = movieCreateDto.Filename
	movie.Duration = &movieCreateDto.Duration
	movie.Thumbnail = &thumb
	movie.ActorID = uint(movieCreateDto.ActorId)
	//movie.Categories = categories
	db.DB.Save(&movie)
	db.DB.Model(&movie).Association("Categories").Clear()
	if len(categories) > 0 {
		db.DB.Model(&movie).Association("Categories").Append(categories)
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/movie/index")

	// entity := models.Movie{
	// 	Title:    movieCreateDto.Title,
	// 	Actor:    &movieCreateDto.Actor,
	// 	Filename: movieCreateDto.Filename,
	// 	Duration: &movieCreateDto.Duration,
	// 	//Categories: categories,
	// 	//Thumbnail: &pictureFileName,
	// }

}
