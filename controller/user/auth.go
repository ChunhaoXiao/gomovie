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

type PasswordDto struct {
	OldPassword          string `form:"oldpassword" binding:"required"`
	Password             string `form:"password" binding:"required,min=6,max=30"`
	PasswordConfirmation string `form:"password_confirmation" binding:"required"`
}

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

func ChangePassword(c *gin.Context) {
	c.HTML(http.StatusOK, "user/auth/changepass.html", gin.H{})
}

func UpdatePassword(c *gin.Context) {
	var passwordDto PasswordDto
	err := c.ShouldBind(&passwordDto)
	if err != nil {
		//fmt.Println("err is#######", err.Error())
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {

			validationErr := utils.GetVfalidationMsg(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": validationErr[0],
			})
			return
		}
	}
	if passwordDto.Password != passwordDto.PasswordConfirmation {
		fmt.Println("p1", passwordDto.Password)
		fmt.Println("p2", passwordDto.PasswordConfirmation)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "两次密码不一致",
		})
		return
	}

	user, _ := c.Get("currentUser")
	userInfo := user.(models.User)
	fmt.Println("pass in database::", userInfo.Password)
	fmt.Println("submitted password", passwordDto.Password)
	passwordErr := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(passwordDto.OldPassword))
	if passwordErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "原密码不正确",
		})
		return
	}
	userInfo.Password = passwordDto.Password
	db.DB.Save(&userInfo)
	c.SetCookie("user", "", 0, "/", "", false, true)
	c.Redirect(http.StatusMovedPermanently, "/auth/login")
}
