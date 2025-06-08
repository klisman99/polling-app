package routes

import (
	"net/http"
	"polling-app/poll-service/handlers"
	"polling-app/poll-service/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(pollService *services.PollService) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTION" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	pollHandler := handlers.NewPollHandler(pollService)

	api := router.Group("/api/v1")
	{
		pollsApi := api.Group("/polls")
		{
			pollsApi.POST("", pollHandler.CreatePoll)
			pollsApi.GET("", pollHandler.GetAllPolls)
			pollsApi.GET("/:id", pollHandler.GetPollByID)
			pollsApi.PUT("/:id", pollHandler.UpdatePoll)
			pollsApi.DELETE("/:id", pollHandler.DeletePoll)
		}
	}

	return router
}
