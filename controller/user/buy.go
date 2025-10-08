package user

import (
	"fmt"
	"movie/db"
	"movie/models"
	"movie/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	Id string `json:"id"`
}

func BuyMovie(c *gin.Context) {
	var movieInfo Movie
	var movie models.Movie
	user, _ := c.Get("currentUser")
	userinfo := user.(models.User)

	err := c.ShouldBind(&movieInfo)

	fmt.Println("movieinfo####", movieInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "参数错误",
		})
		fmt.Println("err", err.Error())
		return
	}
	id := movieInfo.Id
	//mid, _ := strconv.Atoi(id)
	fmt.Println("idddd", id)
	movieErr := db.DB.Where("id=?", id).First(&movie).Error
	fmt.Println("movie####", movie)
	if movieErr != nil {
		fmt.Println("=======================================")
		c.JSON(http.StatusNotFound, gin.H{
			"err": "数据不存在",
		})
		return
	}
	var movies []models.Movie
	db.DB.Model(&userinfo).Where("id=?", id).Association("Movies").Find(&movies)
	if len(movies) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "已经购买过",
		})
		return
	}
	fmt.Println("bought moviesssss###", movies)

	fmt.Println("userinfo.UserCoin.Value#######", userinfo.UserCoin.Value)

	if uint(userinfo.UserCoin.Value) < movie.Price {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"err": "余额不足",
		})
		return
	}

	db.DB.Model(&userinfo).Association("Movies").Append(&movie)

	coin := uint(userinfo.UserCoin.Value) - movie.Price

	db.DB.Model(&userinfo.UserCoin).Update("Value", coin)

	c.JSON(http.StatusOK, gin.H{})

}

func BuyList(c *gin.Context) {
	user, _ := c.Get("currentUser")
	userinfo := user.(models.User)
	var movies []models.Movie
	perPage := 30
	db.DB.Model(&userinfo).Order("id desc").Scopes(utils.Paginate(c.Request, perPage)).Association("Movies").Find(&movies)
	fmt.Println("bought movies:", movies)
	c.HTML(http.StatusOK, "user/buy/index.html", gin.H{
		"movies": movies,
	})

}
