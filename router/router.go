package router

import (
	"fmt"
	"movie/controller"
	"movie/controller/admin"
	"movie/controller/user"
	"movie/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	route := gin.Default()
	//route.Use(ErrorHandler())
	route.Use(middleware.Auth())

	route.Static("/assets", "./assets")
	route.Static("/thumbs", "./thumb")
	route.Static("/actor", "./uploads/actor")
	route.GET("/todo", controller.Home)

	route.GET("/admin/movie/index", admin.MovieList)
	route.GET("/admin/movie/create", admin.Create)
	route.POST("/admin/movie/save", admin.SaveMovie)
	route.GET("/admin/movie/:id", admin.PlayMovie)
	route.DELETE("/admin/movie/:id", admin.RemoveMovie)
	route.GET("/admin/movie/edit/:id", admin.EditMovie)
	//route.POST("/admin/movie/update/:id", admin.UpdateMovie)
	route.GET("/admin/category/create", admin.CreateCategory)
	route.GET("/admin/category/edit/:id", admin.EditCategory)
	route.GET("/admin/category", admin.CategoryList)
	route.POST("/admin/category/update/:id", admin.UpdateCategory)
	route.POST("/admin/category/save", admin.SaveCategory)
	route.DELETE("/admin/category/:id", admin.DeleteCategory)
	route.GET("/admin/actor/create", admin.CreateActor)
	route.POST("/admin/actor/save", admin.SaveActor)
	route.GET("/admin/actor/index", admin.ActorIndex)
	route.GET("/admin/actor/edit/:id", admin.EditActor)
	route.POST("/admin/actor/update/:id", admin.UpdateActor)
	route.DELETE("/admin/actor/:id", admin.DeleteActor)
	route.POST("/admin/movie/checkfile", admin.CheckMovieFile)
	route.GET("/admin/config", admin.ShowConfigForm)
	route.POST("/admin/config/save", admin.SaveConfig)

	route.POST("/admin/movie/upload", admin.Upload)
	route.POST("/admin/actor/upload", admin.UploadPicture)
	route.GET("admin/card/create", admin.CardCreate)
	route.POST("/admin/card/save", admin.SaveCard)
	route.GET("/admin/card/index", admin.GetGroupList)
	route.GET("/admin/card/show/:id", admin.ShowCard)

	route.GET("/index", user.Home)
	route.GET("/video/:id", user.StreamVideo)
	route.GET("/video/list", user.MovieLists)
	route.GET("/video/show/:id", user.Show)

	route.GET("/auth/register", user.RegisterForm)
	route.POST("/auth/register", user.SaveUser)
	route.GET("/auth/login", user.LoginForm)
	route.POST("/auth/login", user.Login)

	//route.GET("/user/charge", user.ChargeForm)

	route.GET("/player/:actorid/movies", user.ActorMovies)
	route.GET("/player/index", user.ActorList)

	grp := route.Group("/user", middleware.RedirectIfNotAuthenticated())
	//grp.Use(middleware.RedirectIfNotAuthenticated())
	grp.GET("/charge", user.ChargeForm)
	grp.POST("/charge", user.DoCharge)
	grp.POST("/buy", user.BuyMovie)
	grp.GET("/buy/index", user.BuyList)
	grp.GET("/auth/changepass", user.ChangePassword)
	grp.POST("/auth/password/update", user.UpdatePassword)
	grp.GET("/auth/logout", user.Logout)
	return route

}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Step1: Process the request first.

		// Step2: Check if any errors were added to the context
		fmt.Println(len(c.Errors))
		if len(c.Errors) > 0 {

			err := c.Errors[0]
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
			// Step3: Use the last error
			// err := c.Errors.Last().Err

			// // Step4: Respond with a generic error message
			// c.JSON(http.StatusInternalServerError, map[string]any{
			// 	"success": false,
			// 	"message": err.Error(),
			// })
		}

		// Any other steps if no errors are found
	}
}
