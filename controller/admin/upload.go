package admin

import (
	"fmt"
	"movie/db"
	"movie/models"
	"movie/utils"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("upload err####", err)
	}
	fmt.Println("file is:###", file)
	fileName := uuid.New().String()
	ext := filepath.Ext(file.Filename)
	uploadedFileName := "./uploads/" + fileName + ext
	c.SaveUploadedFile(file, uploadedFileName)
	thumb := utils.MakeSingleMovieThumb(uploadedFileName)
	c.JSON(200, gin.H{
		"data":  fileName + ext,
		"thumb": thumb,
	})

}

func UploadPicture(c *gin.Context) {
	//form, _ := c.MultipartForm()
	file, _ := c.FormFile("file[]")
	uploadedPictureName := uuid.New().String() + filepath.Ext(file.Filename)
	fmt.Println("uploadedPictureName#########", uploadedPictureName)
	c.SaveUploadedFile(file, "./uploads/actor/"+uploadedPictureName)
	c.JSON(200, gin.H{
		"data":  uploadedPictureName,
		"thumb": "",
	})

	/*
		fileName := uuid.New().String()
		ext := filepath.Ext(file.Filename)
		uploadedFileName := "./uploads/actor/" + fileName + ext
		c.SaveUploadedFile(file, uploadedFileName)
		thumb := utils.MakeMovieThumb(uploadedFileName)
		c.JSON(200, gin.H{
			"data":  fileName + ext,
			"thumb": thumb,
		})*/

}

func CheckMovieFile(c *gin.Context) {
	var data map[string]int
	c.ShouldBind(&data)
	fmt.Println("post data#########", data)
	duration := data["duration"]
	fmt.Println("duration............", duration)
	var movie models.Movie
	err := db.DB.Where("duration=?", duration).First(&movie).Error
	fmt.Println("movie$$$$$$$$$$$$$$$$$$$$$", movie)
	if err != nil {
		c.JSON(200, gin.H{"exist": false})
		return
	}
	c.JSON(200, gin.H{
		"exist": true,
	})
}
