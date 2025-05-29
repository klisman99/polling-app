package server

import (
	"log"
	"os"
	"polling-app/poll-service/internal/db"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
	DB     *db.MongoDB
}

func NewServer() *Server {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}
	mongoDB, err := db.NewMongoDB(mongoURI, "polling_app")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := mongoDB.EnsurePollsCollection(); err != nil {
		log.Fatalf("Failed to ensure polls collection: %v", err)
	}

	r := gin.Default()
	srv := &Server{Router: r, DB: mongoDB}
	srv.setupRoutes()
	return srv
}

func (s *Server) Run(addr string) error {
	defer s.DB.Close()
	return s.Router.Run(addr)
}
