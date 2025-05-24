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

func (s *Server) setupRoutes() {
	s.Router.GET("/health", func(c *gin.Context) {
		if err := s.DB.Client.Ping(c, nil); err != nil {
			c.JSON(500, gin.H{"status": "MongoDB is not reachable"})
			return
		}
		c.JSON(200, gin.H{"status": "Poll Service is running"})
	})
}

func (s *Server) Run(addr string) error {
	defer s.DB.Close()
	return s.Router.Run(addr)
}
