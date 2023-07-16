package routes

import (
	"kebondowo/controllers"
	"kebondowo/middleware"

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
	router.GET("/", middleware.DeserializeUser(), tc.tourismController.GetAll)
	router.GET("/:id", middleware.DeserializeUser(), tc.tourismController.GetOne)
	router.POST("/", tc.tourismController.Create)
	router.PUT("/:id", middleware.DeserializeUser(), tc.tourismController.Update)
	router.DELETE("/:id", middleware.DeserializeUser(), tc.tourismController.Delete)
}