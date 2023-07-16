package main

import (
	"log"
	"net/http"

	"kebondowo/controllers"
	"kebondowo/initializers"
	"kebondowo/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server                 *gin.Engine
	AuthController         controllers.AuthController
	AuthRouteController    routes.AuthRouteController
	UserController         controllers.UserController
	UserRouteController    routes.UserRouteController
	RoleController         controllers.RoleController
	RoleRouteController    routes.RoleRouteController
	TourismController      controllers.TourismController
	TourismRouteController routes.TourismRouteController
	GalleryController      controllers.GalleryController
	GalleryRouteController routes.GalleryRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	RoleController = controllers.NewRoleController(initializers.DB)
	RoleRouteController = routes.NewRoleRouteController(RoleController)

	TourismController = controllers.NewTourismController(initializers.DB)
	TourismRouteController = routes.NewTourismRouteController(TourismController)

	GalleryController = controllers.NewGalleryController(initializers.DB)
	GalleryRouteController = routes.NewGalleryRouteController(GalleryController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	RoleRouteController.RoleRoute(router)
	TourismRouteController.TourismRoute(router)
	GalleryRouteController.GalleryRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
