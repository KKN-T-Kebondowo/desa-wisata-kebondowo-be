package routes

import (
	"kebondowo/controllers"
	"kebondowo/middleware"

	"github.com/gin-gonic/gin"
)

type ArticleRouteController struct {
	articleController controllers.ArticleController
}

func NewArticleRouteController(articleController controllers.ArticleController) ArticleRouteController {
	return ArticleRouteController{articleController}
}

func (ac *ArticleRouteController) ArticleRoute(rg *gin.RouterGroup) {
	router := rg.Group("/articles")
	router.GET("/", ac.articleController.GetAll)
	router.GET("/:slug", ac.articleController.GetOne) // Use :slug instead of :id
	router.POST("/", middleware.DeserializeUser() ,ac.articleController.Create)
	router.PUT("/:id", middleware.DeserializeUser(),ac.articleController.Update)
	router.DELETE("/:id", middleware.DeserializeUser() ,ac.articleController.Delete)
}
