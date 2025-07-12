package admin

import (
	"movie/db"
	"movie/dto"
	"movie/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/category/create.html", gin.H{})
}

func SaveCategory(c *gin.Context) {
	var category models.Category
	c.Bind(&category)
	db.DB.Create(&category)
	c.Redirect(http.StatusMovedPermanently, "/admin/category")

}

func CategoryList(c *gin.Context) {
	var categories []models.Category
	db.DB.Find(&categories)
	c.HTML(http.StatusOK, "admin/category/index.html", gin.H{
		"datas": categories,
	})
}

func EditCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	err := db.DB.Where("id=?", id).First(&category).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusBadRequest,
			"message": "数据不存在",
		})
		return
	}
	c.HTML(http.StatusOK, "admin/category/create.html", gin.H{
		"name": category.Name,
		"id":   category.ID,
	})

}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	err := db.DB.First(&category, id).Error
	if err != nil {
		return
	}
	var categoryDto dto.CategoryDto
	c.Bind(&categoryDto)
	category.Name = categoryDto.Name
	db.DB.Save(&category)
	c.Redirect(http.StatusMovedPermanently, "/admin/category")

	//category.ID, _ := strconv.Atoi(id)

}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	db.DB.Where("id=?", id).Delete(&models.Category{})
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusBadRequest,
		"message": "数据不存在",
	})
}
