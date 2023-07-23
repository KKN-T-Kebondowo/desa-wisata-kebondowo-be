package routes

import (
	"kebondowo/controllers"

	"github.com/gin-gonic/gin"
)

type DashboardRouteController struct {
	dashboardController controllers.DashboardController
}

func NewDashboardRouteController(dashboardController controllers.DashboardController) DashboardRouteController {
	return DashboardRouteController{dashboardController}
}

func (dc *DashboardRouteController) DashboardRoute(rg *gin.RouterGroup) {
	router := rg.Group("/dashboard")

	router.GET("/", dc.dashboardController.GetDashboard)
}