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

	initializers.DB.AutoMigrate(&models.Role{})
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Tourism{})
	initializers.DB.AutoMigrate(&models.TourismPicture{})
	initializers.DB.AutoMigrate(&models.Gallery{})
	initializers.DB.AutoMigrate(&models.Article{})
	initializers.DB.AutoMigrate(&models.UMKM{})
	initializers.DB.AutoMigrate(&models.UMKMPicture{})

	fmt.Println("? Migration complete")
}
