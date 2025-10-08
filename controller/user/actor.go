package user

import (
	"fmt"
	"movie/db"
	"movie/models"
	"movie/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ActorList(c *gin.Context) {
	var actorList []models.Actor
	subQuery := db.DB.Model(&models.Movie{}).Select("actor_id", "COUNT(id) as movie_count").Group("actor_id")
	db.DB.Joins("LEFT JOIN (?) as counts on actors.id=counts.actor_id", subQuery).Order("counts.movie_count DESC").Find(&actorList)
	fmt.Println("actor lists:::", actorList)
	c.HTML(http.StatusOK, "user/actor/lists.html", gin.H{
		"actors": actorList,
	})
}

func ActorMovies(c *gin.Context) {
	actorId := c.Param("actorid")
	var actor models.Actor
	err := db.DB.Where("id=?", actorId).First(&actor).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "数据不存在!!!!!!",
		})
	}
	fmt.Println("actor:#######", actor)
	var movies []models.Movie
	perPage := 20

	db.DB.Model(&actor).Scopes(utils.Paginate(c.Request, perPage)).Association("Movies").Find(&movies)
	var movieCount int64
	db.DB.Model(&models.Movie{}).Where("actor_id=?", actor.ID).Count(&movieCount)

	currentPage, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageArr := utils.GeneratePageStr(perPage, currentPage, int(movieCount), 5)
	fmt.Print(pageArr)
	c.HTML(http.StatusOK, "user/actor/movies.html", gin.H{
		"actor":       actor,
		"movies":      movies,
		"pageArr":     pageArr,
		"currentPage": currentPage,
	})

}
