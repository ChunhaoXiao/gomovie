package middleware

import (
	"movie/dto"

	"github.com/gin-gonic/gin"
)

var user = &dto.LoginDto{
	Username: "",
}

func GetUserName() string {
	return user.Username
}

func Auth() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Pragma", "No-cache")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		cvalue, error := c.Cookie("user")
		if error != nil {
			user.Username = ""
			//c.AddToMap(gin.H{})
			//c.AddToMap(gin.H{"loginUser": ""})
			//c.Redirect(http.StatusMovedPermanently, "/user/login")
		} else {
			user.Username = cvalue

		}

		c.Next()

	}
}
