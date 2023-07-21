package routes

import (
	"kebondowo/controllers"

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
	router.GET("/:id", ac.articleController.GetOne)
	router.POST("/", ac.articleController.Create)
	router.PUT("/:id", ac.articleController.Update)
	router.DELETE("/:id", ac.articleController.Delete)
}
