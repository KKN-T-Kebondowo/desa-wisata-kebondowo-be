package routes

import (
	"kebondowo/controllers"

	"github.com/gin-gonic/gin"
)

type TourismRouteController struct {
	tourismController controllers.TourismController
}

func NewTourismRouteController(tourismController controllers.TourismController) TourismRouteController {
	return TourismRouteController{tourismController}
}

func (tc *TourismRouteController) TourismRoute(rg *gin.RouterGroup) {
	
	router := rg.Group("tourisms")
	router.GET("/",  tc.tourismController.GetAll)
	router.GET("/:id", tc.tourismController.GetOne)
	router.POST("/", tc.tourismController.Create)
	router.PUT("/:id",  tc.tourismController.Update)
	router.DELETE("/:id", tc.tourismController.Delete)
}