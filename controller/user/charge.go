package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChargeForm(c *gin.Context) {

	c.HTML(http.StatusOK, "user/charge/create.html", gin.H{})
}

type ChargeData struct {
	CardNumber string
}

func DoCharge(c *gin.Context) {
	var number ChargeData
	c.ShouldBind(&number)

}
