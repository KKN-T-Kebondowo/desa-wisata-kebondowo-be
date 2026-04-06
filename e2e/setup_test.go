package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"kebondowo/controllers"
	"kebondowo/initializers"
	"kebondowo/models"
	"kebondowo/routes"
	"kebondowo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testServer *httptest.Server
	testDB     *gorm.DB
	authToken  string
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// Load config from project root
	os.Chdir("..")
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to test database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DBHost, config.DBUserName, config.DBUserPassword, "kebondowo_test", config.DBPort,
	)
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	// Set global DB for middleware
	initializers.DB = testDB

	// Migrate tables
	migrationModels := []interface{}{
		&models.Role{},
		&models.User{},
		&models.Tourism{},
		&models.TourismPicture{},
		&models.Gallery{},
		&models.Article{},
		&models.UMKM{},
		&models.UMKMPicture{},
	}
	for _, model := range migrationModels {
		if err := testDB.AutoMigrate(model); err != nil {
			log.Fatalf("failed to migrate: %v", err)
		}
	}

	// Clean all tables
	cleanDB()

	// Seed roles
	testDB.Create(&models.Role{Name: "admin"})
	testDB.Create(&models.Role{Name: "user"})

	// Setup server
	server := setupServer(testDB)
	testServer = httptest.NewServer(server)
	defer testServer.Close()

	// Seed admin user and get auth token
	seedAdminAndLogin(config)

	// Run tests
	code := m.Run()

	// Cleanup
	cleanDB()
	os.Exit(code)
}

func cleanDB() {
	testDB.Exec("DELETE FROM umkm_pictures")
	testDB.Exec("DELETE FROM umkms")
	testDB.Exec("DELETE FROM tourism_pictures")
	testDB.Exec("DELETE FROM tourisms")
	testDB.Exec("DELETE FROM articles")
	testDB.Exec("DELETE FROM galleries")
	testDB.Exec("DELETE FROM users")
	testDB.Exec("DELETE FROM roles")
}

func setupServer(db *gorm.DB) *gin.Engine {
	server := gin.New()

	authController := controllers.NewAuthController(db)
	userController := controllers.NewUserController(db)
	roleController := controllers.NewRoleController(db)
	tourismController := controllers.NewTourismController(db)
	galleryController := controllers.NewGalleryController(db)
	articleController := controllers.NewArticleController(db)
	tourismPictureController := controllers.NewTourismPictureController(db)
	dashboardController := controllers.NewDashboardController(db)
	umkmController := controllers.NewUMKMController(db)
	umkmPictureController := controllers.NewUMKMPictureController(db)

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

	router := server.Group("/api")

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

	return server
}

func seedAdminAndLogin(config initializers.Config) {
	// Get the admin role ID
	var adminRole models.Role
	testDB.Where("name = ?", "admin").First(&adminRole)

	// Create admin user via direct DB insert (since register now requires auth)
	hashedPw, _ := hashForTest("admin1234")
	testDB.Create(&models.User{
		Username: "admin",
		Password: hashedPw,
		RoleID:   adminRole.ID,
	})

	// Login to get token
	body := map[string]string{
		"username": "admin",
		"password": "admin1234",
	}
	resp := doRequest("POST", "/api/auth/login", body, "")
	var result map[string]interface{}
	json.Unmarshal(resp.Body, &result)
	if token, ok := result["access_token"].(string); ok {
		authToken = token
	} else {
		log.Fatalf("failed to get auth token from login response: %v", result)
	}
}

func hashForTest(password string) (string, error) {
	return utils.HashPassword(password)
}

// Response helpers
type testResponse struct {
	StatusCode int
	Body       []byte
}

func doRequest(method, path string, body interface{}, token string) testResponse {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, _ := http.NewRequest(method, testServer.URL+path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return testResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
	}
}

func parseJSON(body []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	return result
}
