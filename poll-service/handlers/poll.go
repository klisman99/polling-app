package handlers

import (
	"net/http"
	"polling-app/poll-service/models"
	"polling-app/poll-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PollHandler struct {
	service   *services.PollService
	validator *validator.Validate
}

func NewPollHandler(pollService *services.PollService) *PollHandler {
	return &PollHandler{
		service:   pollService,
		validator: validator.New(),
	}
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func validateNotFoundAndRespond(c *gin.Context, err error, message string) bool {
	if err == nil {
		return true
	}

	if err.Error() == "poll not found" {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Poll not found",
		})
	} else {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: message,
		})
	}

	return false
}

func (h *PollHandler) CreatePoll(c *gin.Context) {
	var req models.CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validation failed",
			Data:    err.Error(),
		})
		return
	}

	poll, err := h.service.CreatePoll(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to create poll",
		})
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Poll created successfully",
		Data:    poll,
	})
}

func (h *PollHandler) GetAllPolls(c *gin.Context) {
	polls, err := h.service.GetAllPolls(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to fetch polls",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Polls fetched successfully",
		Data:    polls,
	})
}

func (h *PollHandler) GetPollByID(c *gin.Context) {
	id := c.Param("id")

	poll, err := h.service.GetPollById(c, id)

	if !validateNotFoundAndRespond(c, err, "Failed to fetch poll") {
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Fetched poll successfully",
		Data:    poll,
	})
}

func (h *PollHandler) UpdatePoll(c *gin.Context) {
	id := c.Param("id")

	var req *models.UpdatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Validation failed",
		})
		return
	}

	poll, err := h.service.UpdatePoll(c, id, req)
	if !validateNotFoundAndRespond(c, err, "Failed to update poll") {
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Poll updated successfully",
		Data:    poll,
	})
}

func (h *PollHandler) DeletePoll(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeletePoll(c, id)
	if !validateNotFoundAndRespond(c, err, "Failed to delete poll") {
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Poll deleted successfully",
	})
}
