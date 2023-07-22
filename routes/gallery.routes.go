package routes

import (
	"kebondowo/controllers"

	"github.com/gin-gonic/gin"
)

type GalleryRouteController struct {
	galleryController controllers.GalleryController
}

func NewGalleryRouteController(galleryController controllers.GalleryController) GalleryRouteController {
	return GalleryRouteController{galleryController}
}

func (gc *GalleryRouteController) GalleryRoute(rg *gin.RouterGroup) {
	router := rg.Group("/galleries")
	router.GET("/", gc.galleryController.GetAll)
	router.GET("/:id", gc.galleryController.GetOne)
	router.POST("/", gc.galleryController.Create)
	router.DELETE("/:id", gc.galleryController.Delete)
}
