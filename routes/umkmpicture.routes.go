package routes

import (
	"kebondowo/controllers"
	"kebondowo/middleware"

	"github.com/gin-gonic/gin"
)

type UMKMPictureRouteController struct {

	umkmPictureController controllers.UMKMPictureController
}

func NewUMKMPictureRouteController(umkmPictureController controllers.UMKMPictureController) UMKMPictureRouteController {
	return UMKMPictureRouteController{umkmPictureController}
}

func (upc *UMKMPictureRouteController) UMKMPictureRoute(rg *gin.RouterGroup) {
	
	router := rg.Group("/umkm-pictures")
	router.GET("/", upc.umkmPictureController.GetAll)
	router.GET("/:id", upc.umkmPictureController.GetOne)
	router.POST("/", middleware.DeserializeUser(), upc.umkmPictureController.Create)
	router.DELETE("/:id", middleware.DeserializeUser(), upc.umkmPictureController.Delete)
}

