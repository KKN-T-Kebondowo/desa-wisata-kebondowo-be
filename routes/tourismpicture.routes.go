package routes

import (
	"kebondowo/controllers"

	"github.com/gin-gonic/gin"
)

type TourismPictureRouteController struct {
	tourismPictureController controllers.TourismPictureController
}

func NewTourismPictureRouteController(tourismPictureController controllers.TourismPictureController) TourismPictureRouteController {
	return TourismPictureRouteController{tourismPictureController}
}

func (tpc *TourismPictureRouteController) TourismPictureRoute(rg *gin.RouterGroup) {
	router := rg.Group("/tourism-pictures")
	router.GET("/", tpc.tourismPictureController.GetAll)
	router.GET("/:id", tpc.tourismPictureController.GetOne)
	router.POST("/", tpc.tourismPictureController.Create)
	router.DELETE("/:id", tpc.tourismPictureController.Delete)
}


