package user

import (
	"fmt"
	"math"
	"movie/db"
	"movie/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ChargeForm(c *gin.Context) {
	// username, _ := c.Cookie("user")
	// fmt.Println("user#######", username)
	// var user models.User
	// db.DB.Where("username=?", username).First(&user)
	// var usercoin models.UserCoin
	// db.DB.Where("user_id=?").First(&usercoin)
	user, _ := c.Get("currentUser")
	info := user.(models.User)
	fmt.Println("info ---->", info.UserCoin)
	c.HTML(http.StatusOK, "user/charge/create.html", gin.H{"coin": info.UserCoin.Value})
}

type ChargeData struct {
	CardNumber string `json:"card_number"`
}

func DoCharge(c *gin.Context) {
	_, err := c.Cookie("user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户不存在",
		})
		return
	}
	username, _ := c.Cookie("user")
	var user models.User
	db.DB.Where("username=?", username).First(&user)
	fmt.Println("user==============>", user)
	var limit models.ChargeLimit
	limiterr := db.DB.Where("user_id=?", user.ID).First(&limit).Error
	if limiterr == nil {
		s := math.Floor(time.Since(limit.UpdatedAt).Seconds())
		fmt.Println("secondss", s)
		if s < 300 {
			if limit.TryTimes > 4 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "错误次数太多，请稍后再试",
				})
				return
			}
		} else {
			fmt.Println("deleting...........")
			db.DB.Unscoped().Delete(&limit)
		}

	}

	var number ChargeData
	c.ShouldBind(&number)
	cardNumber := number.CardNumber
	fmt.Println("card number:", cardNumber)
	if cardNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "卡号不能为空",
		})
		return
	}
	var card models.Card

	error := db.DB.Scopes(models.UnusedCard(cardNumber)).First(&card).Error
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "卡号不正确",
		})
		var chargeLimit models.ChargeLimit
		err := db.DB.Where("user_id=?", user.ID).First(&chargeLimit).Error
		if err != nil {
			limit := models.ChargeLimit{
				UserId:   int(user.ID),
				TryTimes: 1,
			}
			db.DB.Create(&limit)
		} else {
			chargeLimit.TryTimes = chargeLimit.TryTimes + 1
			db.DB.Save(&chargeLimit)
		}
		return
	}
	var usercoin models.UserCoin

	errors := db.DB.Where("user_id=?", user.ID).First(&usercoin).Error
	if errors != nil {
		data := models.UserCoin{
			UserID: int(user.ID),
			Value:  card.CoinValue,
		}
		db.DB.Create(&data)
	} else {
		value := usercoin.Value + card.CoinValue
		db.DB.Model(&usercoin).Update("Value", value)
	}
	db.DB.Model(&card).Update("user_id", user.ID)
	db.DB.Where("user_id=?", user.ID).Delete(&models.ChargeLimit{})
	c.JSON(http.StatusOK, gin.H{
		"data": "success",
	})

}
