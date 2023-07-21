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
	ArticleController      controllers.ArticleController
	ArticleRouteController routes.ArticleRouteController
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

	ArticleController = controllers.NewArticleController(initializers.DB)
	ArticleRouteController = routes.NewArticleRouteController(ArticleController)

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
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"} // Add "Content-Type" to allowed headers

	server := gin.Default()
	server.Use(cors.New(corsConfig))

	// server := gin.Default()

	// Enable CORS
	// server.Use(CORS())

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	RoleRouteController.RoleRoute(router)
	TourismRouteController.TourismRoute(router)
	GalleryRouteController.GalleryRoute(router)
	ArticleRouteController.ArticleRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}

// CORS middleware
// func CORS() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusNoContent)
// 			return
// 		}

// 		c.Next()
// 	}
// }
