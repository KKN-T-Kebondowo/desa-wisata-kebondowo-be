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

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables: ", err)
	}

	initializers.ConnectDB(&config)

	// Initialize controllers
	authController := controllers.NewAuthController(initializers.DB)
	userController := controllers.NewUserController(initializers.DB)
	roleController := controllers.NewRoleController(initializers.DB)
	tourismController := controllers.NewTourismController(initializers.DB)
	galleryController := controllers.NewGalleryController(initializers.DB)
	articleController := controllers.NewArticleController(initializers.DB)
	tourismPictureController := controllers.NewTourismPictureController(initializers.DB)
	dashboardController := controllers.NewDashboardController(initializers.DB)
	umkmController := controllers.NewUMKMController(initializers.DB)
	umkmPictureController := controllers.NewUMKMPictureController(initializers.DB)

	// Initialize route controllers
	authRouteController := routes.NewAuthRouteController(authController)
	userRouteController := routes.NewRouteUserController(userController)
	roleRouteController := routes.NewRoleRouteController(roleController)
	tourismRouteController := routes.NewTourismRouteController(tourismController)
	galleryRouteController := routes.NewGalleryRouteController(galleryController)
	articleRouteController := routes.NewArticleRouteController(articleController)
	tourismPictureRouteController := routes.NewTourismPictureRouteController(tourismPictureController)
	dashboardRouteController := routes.NewDashboardRouteController(dashboardController)
	umkmRouteController := routes.NewUMKMRouteController(umkmController)
	umkmPictureRouteController := routes.NewUMKMPictureRouteController(umkmPictureController)

	// Setup server
	server := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "https://dashboard-desa-wisata-kebondowo.vercel.app", "https://desa-wisata-kebondowo-dashboard.vercel.app", "https://dashboard.kebondowo.com"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	authRouteController.AuthRoute(router)
	userRouteController.UserRoute(router)
	roleRouteController.RoleRoute(router)
	tourismRouteController.TourismRoute(router)
	galleryRouteController.GalleryRoute(router)
	articleRouteController.ArticleRoute(router)
	tourismPictureRouteController.TourismPictureRoute(router)
	dashboardRouteController.DashboardRoute(router)
	umkmRouteController.UMKMRoute(router)
	umkmPictureRouteController.UMKMPictureRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
