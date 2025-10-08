package main

import (
	"movie/db"
	"movie/middleware"
	"movie/models"
	"movie/router"
	"movie/templatefunc"
	"movie/utils"
	"text/template"
)

func main() {

	category := models.Category{}
	movie := models.Movie{}
	actor := models.Actor{}
	user := models.User{}
	cardgroup := models.Cardgroup{}
	card := models.Card{}
	userCoin := models.UserCoin{}
	ChargeLimit := models.ChargeLimit{}
	Configration := models.Configuration{}
	//userBuy := models.UserBuy{}
	db.ConnectDB()
	//db.DB.AutoMigrate(&todo)
	db.DB.AutoMigrate(&category)
	db.DB.AutoMigrate(&movie)
	db.DB.AutoMigrate(&actor)
	db.DB.AutoMigrate(&user)
	db.DB.AutoMigrate(&cardgroup)
	db.DB.AutoMigrate(&card)
	db.DB.AutoMigrate(&userCoin)
	db.DB.AutoMigrate(&ChargeLimit)
	db.DB.AutoMigrate(&Configration)
	//db.DB.AutoMigrate(&userBuy)
	r := router.InitRoute()

	r.SetFuncMap(template.FuncMap{
		"formatDuration":  templatefunc.FormatDuration,
		"getActors":       templatefunc.GetActors,
		"getCategories":   templatefunc.GetCategories,
		"actorName":       templatefunc.ActorName,
		"indexActors":     templatefunc.UserIndexActor,
		"loginUser":       middleware.GetUserName,
		"getChargeConfig": templatefunc.GetChargeConfig,
		//"getThumbNail":   templatefunc.GetThunmbNail,
		//"selectedCategory": templatefunc.IsCategorySelected,
	})
	utils.SeedMovie()
	r.LoadHTMLGlob("template/**/**/*")
	r.Run(":8081")
}
