package server

import "github.com/gin-gonic/gin"

type Server struct {
	Router *gin.Engine
}

func NewServer() *Server {
	r := gin.Default()
	srv := &Server{Router: r}
	srv.setupRoutes()
	return srv
}

func (s *Server) setupRoutes() {
	s.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Poll Service is running"})
	})
}

func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}
