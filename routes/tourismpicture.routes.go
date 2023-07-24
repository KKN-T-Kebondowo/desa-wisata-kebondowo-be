package routes

import (
	"kebondowo/controllers"
	"kebondowo/middleware"

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
	router.POST("/", middleware.DeserializeUser(),tpc.tourismPictureController.Create)
	router.DELETE("/:id", middleware.DeserializeUser(),tpc.tourismPictureController.Delete)
}


