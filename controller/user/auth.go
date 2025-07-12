package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"movie/db"
	"movie/dto"
	"movie/models"
	"movie/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func RegisterForm(c *gin.Context) {
	c.HTML(http.StatusOK, "user/auth/register.html", gin.H{})
}

func SaveUser(c *gin.Context) {
	var userCreate dto.UserCreateDto
	if err := c.ShouldBindWith(&userCreate, binding.JSON); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {

			validationErr := utils.GetVfalidationMsg(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": validationErr[0],
			})
		} else {
			fmt.Println("Non-validator error:", err)
		}
		return
	}
	if userCreate.Password != userCreate.RePassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": "两次密码不一致",
		})
		return
	}
	var user models.User
	err := db.DB.Where("username=?", userCreate.Username).First(&user).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": "用户名已存在",
		})
		return
	}
	emailErr := db.DB.Where("email=?", userCreate.Email).First(&user).Error
	if emailErr == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": "邮箱已被占用",
		})
		return
	}
	createUser := models.User{
		Username: userCreate.Username,
		Email:    userCreate.Email,
		Password: userCreate.Password,
	}
	db.DB.Create(&createUser)

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func LoginForm(c *gin.Context) {
	referer := c.GetHeader("Referer")
	fmt.Println("refere#######", referer)
	if strings.Contains(referer, "/show") {
		c.SetCookie("refer", referer, 60, "/", "", false, true)
	}
	c.HTML(http.StatusOK, "user/auth/login.html", gin.H{})
}

func Login(c *gin.Context) {
	var loginDto dto.LoginDto
	fmt.Println("login......")
	if err := c.ShouldBindWith(&loginDto, binding.Form); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			validationErr := utils.GetVfalidationMsg(validationErrors)
			fmt.Println("LOGIN VALIDATION erroS:", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": validationErr[0],
			})
		}

		return
	}
	var user models.User
	userErr := db.DB.Where("username=?", loginDto.Username).First(&user).Error
	if userErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": "用户名或密码错误",
		})
		return
	}
	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))
	if passwordErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": "用户名或密码错误！",
		})
		return
	}
	c.SetCookie("user", user.Username, 3600*24, "/", "", false, true)
	refer, err := c.Cookie("refer")
	if err == nil {
		c.Redirect(http.StatusMovedPermanently, refer)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/index")
}

func Logout(c *gin.Context) {
	c.SetCookie("user", "", 0, "/", "", false, true)

	c.Redirect(http.StatusMovedPermanently, "/index")

}
