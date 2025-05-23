package main

import (
	"log"
	"polling-app/poll-service/internal/server"
)

func main() {
	srv := server.NewServer()
	if err := srv.Run(":3002"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
