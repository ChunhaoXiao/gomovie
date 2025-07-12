package admin

import (
	"fmt"
	"movie/utils"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	fmt.Println(file)
	fileName := uuid.New().String()
	ext := filepath.Ext(file.Filename)
	uploadedFileName := "./uploads/" + fileName + ext
	c.SaveUploadedFile(file, uploadedFileName)
	thumb := utils.MakeMovieThumb(uploadedFileName)
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
