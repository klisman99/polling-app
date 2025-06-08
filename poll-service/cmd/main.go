package main

import (
	"log"
	"os"
	"polling-app/poll-service/config"
	"polling-app/poll-service/repositories"
	"polling-app/poll-service/routes"
	"polling-app/poll-service/services"
)

func main() {
	db := config.NewDatabase()
	defer db.Close()

	pollRepo := repositories.NewPollRepository(db.DB)
	pollService := services.NewPollService(pollRepo)
	router := routes.SetupRoutes(pollService)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is no set")
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
	log.Print("Server is running")
}
