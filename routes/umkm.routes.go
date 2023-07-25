package routes

import (
	"kebondowo/controllers"
	"kebondowo/middleware"

	"github.com/gin-gonic/gin"
)

type UMKMRouteController struct {
	umkmController controllers.UMKMController
}

func NewUMKMRouteController(umkmController controllers.UMKMController) UMKMRouteController {
	return UMKMRouteController{umkmController}
}

func (uc *UMKMRouteController) UMKMRoute(rg *gin.RouterGroup) {
	
	router := rg.Group("umkms")
	router.GET("/", uc.umkmController.GetAll)
	router.GET("/:slug", uc.umkmController.GetOne)
	router.POST("/", middleware.DeserializeUser(), uc.umkmController.Create)
	router.PUT("/:id", middleware.DeserializeUser(), uc.umkmController.Update)
	router.DELETE("/:id", middleware.DeserializeUser(), uc.umkmController.Delete)
}

