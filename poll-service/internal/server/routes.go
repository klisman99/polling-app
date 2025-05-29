package server

import (
	"context"
	"net/http"
	"polling-app/poll-service/internal/models"
	"polling-app/poll-service/internal/poll"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) setupRoutes() {
	pollService := poll.NewService(s.DB)

	s.Router.GET("/health", func(c *gin.Context) {
		if err := s.DB.Client.Ping(c, nil); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "MongoDB connection failed"})
		}
		c.JSON(http.StatusOK, gin.H{"status": "Poll service running"})
	})

	polls := s.Router.Group("polls")
	{
		polls.POST("", func(c *gin.Context) {
			var poll models.Poll
			if err := c.ShouldBindJSON(&poll); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid poll data"})
				return
			}

			ctx, cancel := context.WithTimeout(c, 5*time.Second)
			defer cancel()

			id, err := pollService.Create(ctx, &poll)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"id": id})
		})

		polls.GET("", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(c, 5*time.Second)
			defer cancel()

			polls, err := pollService.GetAll(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch polls"})
				return
			}
			c.JSON(http.StatusOK, polls)
		})

		polls.GET("/:id", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(c, 5*time.Second)
			defer cancel()

			id := c.Param("id")
			poll, err := pollService.GetByID(ctx, id)

			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
				return
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find poll"})
				return
			}

			c.JSON(http.StatusOK, poll)
		})

		polls.PUT("/:id", func(c *gin.Context) {
			var poll models.Poll
			if err := c.ShouldBindJSON(&poll); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid poll data"})
				return
			}

			ctx, cancel := context.WithTimeout(c, 5*time.Second)
			defer cancel()

			id := c.Param("id")
			err := pollService.Update(ctx, id, &poll)

			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
				return
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update poll"})
				return
			}

			c.Status(http.StatusOK)
		})

		polls.DELETE("/:id", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(c, 5*time.Second)
			defer cancel()

			id := c.Param("id")
			err := pollService.Delete(ctx, id)

			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
				return
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete poll"})
				return
			}

			c.Status(http.StatusOK)
		})
	}
}
