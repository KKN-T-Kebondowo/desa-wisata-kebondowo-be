package main

import (
	"fmt"
	"kebondowo/initializers"
	"kebondowo/models"
	"log"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
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
		if err := initializers.DB.AutoMigrate(model); err != nil {
			log.Fatalf("failed to migrate %T: %v", model, err)
		}
	}

	fmt.Println("Migration complete")
}
