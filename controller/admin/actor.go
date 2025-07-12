package admin

import (
	"fmt"
	"movie/db"
	"movie/dto"
	"movie/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateActor(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/actor/create.html", gin.H{})
}

func ActorIndex(c *gin.Context) {
	var actors []models.Actor
	db.DB.Find(&actors)
	fmt.Println(actors)
	c.HTML(http.StatusOK, "admin/actor/index.html", gin.H{
		"datas": actors,
	})
}

func SaveActor(c *gin.Context) {
	var actor dto.ActorCreateDto
	c.ShouldBind(&actor)
	fmt.Println(actor)
	entity := models.Actor{
		Name:        actor.Name,
		Pictures:    actor.Pictures,
		IsRecommand: actor.IsRecommand,
	}
	db.DB.Create(&entity)
	c.Redirect(http.StatusMovedPermanently, "/admin/actor/index")

}

func EditActor(c *gin.Context) {
	id := c.Param("id")
	var actor models.Actor
	db.DB.Where("id=?", id).First(&actor)
	c.HTML(http.StatusOK, "admin/actor/create.html", gin.H{
		"actor": actor,
	})
}

func UpdateActor(c *gin.Context) {
	var actorDto dto.ActorCreateDto
	c.ShouldBind(&actorDto)
	fmt.Println("pictures#######", actorDto.Pictures)

	id := c.Param("id")
	var actor models.Actor
	db.DB.Where("id=?", id).First(&actor)
	actor.Name = actorDto.Name
	actor.Pictures = actorDto.Pictures
	actor.IsRecommand = actorDto.IsRecommand
	db.DB.Save(actor)
	c.Redirect(http.StatusMovedPermanently, "/admin/actor/index")
}

func DeleteActor(c *gin.Context) {
	id := c.Param("id")
	var actor models.Actor
	db.DB.First(&actor, id)
	db.DB.Delete(&actor)
	c.JSON(http.StatusOK, gin.H{
		"data": "success",
	})
	//c.Redirect(http.StatusMovedPermanently, "/admin/actor/index")
}
