package routes

import (
	"kebondowo/controllers"
	"kebondowo/middleware"

	"github.com/gin-gonic/gin"
)

type RoleRouteController struct {
	roleController controllers.RoleController
}

func NewRoleRouteController(roleController controllers.RoleController) RoleRouteController {
	return RoleRouteController{roleController}
}

func (rc *RoleRouteController) RoleRoute(rg *gin.RouterGroup) {
	
	router := rg.Group("roles")
	router.GET("/", middleware.DeserializeUser(), rc.roleController.GetAll)
	router.GET("/:id", middleware.DeserializeUser(), rc.roleController.GetOne)
	router.POST("/", rc.roleController.Create)
	router.PUT("/:id", middleware.DeserializeUser(), rc.roleController.Update)
	router.DELETE("/:id", middleware.DeserializeUser(), rc.roleController.Delete)
}
