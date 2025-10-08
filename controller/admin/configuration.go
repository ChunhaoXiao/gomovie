package admin

import (
	"movie/db"
	"movie/dto"
	"movie/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowConfigForm(c *gin.Context) {
	var configs []models.Configuration
	db.DB.Find(&configs)

	datas := make(map[string]string)
	for _, item := range configs {
		datas[item.KeyName] = item.Value
	}
	c.HTML(http.StatusOK, "admin/config/create.html", gin.H{
		"configs": datas,
	})
}

func SaveConfig(c *gin.Context) {
	var allConfigs []models.Configuration
	db.DB.Find(&allConfigs)
	result := make(map[string]models.Configuration)
	for _, item := range allConfigs {
		result[item.KeyName] = item
	}
	var data dto.Configuration
	c.ShouldBind(&data)
	config := models.Configuration{}
	if data.Charge != "" {
		v, ok := result["charge"]
		if ok {
			v.Value = data.Charge
			db.DB.Save(&v)
		} else {
			config.KeyName = "charge"
			config.Value = data.Charge
			db.DB.Create(&config)
		}

	}

	//db.DB.Create(&config)

}
