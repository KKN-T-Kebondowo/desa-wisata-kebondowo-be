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
	initializers.DB.Migrator().CreateConstraint(&models.User{}, "Role")
	initializers.DB.AutoMigrate(&models.User{})

	
	fmt.Println("? Migration complete")
}
