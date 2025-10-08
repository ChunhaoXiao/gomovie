package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RedirectIfNotAuthenticated() gin.HandlerFunc {
	fmt.Println("......................................................")
	return func(ctx *gin.Context) {
		fmt.Println("......................................................###")
		cookieVal, err := ctx.Cookie("user")
		fmt.Println("err is:", err)
		fmt.Println("cookieVal is:", cookieVal)
		if err != nil || cookieVal == "" {
			ctx.Redirect(http.StatusMovedPermanently, "/auth/login ")

		}
		ctx.Next()
	}
}
