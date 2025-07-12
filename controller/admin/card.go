package admin

import (
	"fmt"
	"net/http"
	"time"

	"movie/db"
	"movie/dto"
	"movie/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupCards struct {
	Id         int
	Group_name string
	Card_count int
}

func CardCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/card/create.html", gin.H{})
}

func GetGroupList(c *gin.Context) {
	//var cardGroups []models.Cardgroup
	var cg []GroupCards

	//db.DB.Order("id desc").Find(&cardGroups)
	db.DB.Raw("SELECT cg.*, count(card.id) as card_count FROM cardgroups cg left join cards card on card.cardgroup_id=cg.id group by cg.id order by cg.id desc").Scan(&cg)
	fmt.Println(cg)
	//Raw("SELECT p.*, count(i.id) as invoice_count FROM profiles p left join invoices i on i.profile_fk = p.id group by p.id").
	//Scan(&result)
	c.HTML(http.StatusOK, "admin/card/index.html", gin.H{
		"datas": cg,
	})
}

func ShowCard(c *gin.Context) {
	id := c.Param("id")
	var cards []models.Card
	db.DB.Where("cardgroup_id=?", id).Find(&cards)
	c.HTML(http.StatusOK, "admin/card/show.html", gin.H{
		"datas": cards,
	})
}

func SaveCard(c *gin.Context) {
	var cardDto dto.CardDto
	err := c.ShouldBind(&cardDto)
	fmt.Println("$$$$$$$", cardDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	if cardDto.Quantity > 0 {
		fmt.Println("quantity.......", cardDto.Quantity)
		if cardDto.Coin > 0 {
			fmt.Println("coin.......", cardDto.Coin)
			cardList := []models.Card{}
			for i := 0; i < cardDto.Quantity; i++ {
				id := uuid.New()

				card := models.Card{
					CardNumber: id.String(),
					CoinValue:  cardDto.Coin,
				}
				cardList = append(cardList, card)

			}
			cardGroup := models.Cardgroup{
				GroupName: time.Now().Format("20060102150405"),
				Cards:     cardList,
			}
			fmt.Println("card group$$$$$$$$$$$", cardGroup)
			db.DB.Create(&cardGroup)

		}
	}
}
