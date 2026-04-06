package handlers

import (
	"net/http"
	"time"

	"simple/models"

	"github.com/gin-gonic/gin"
)

// HelloHandler handles the guaranteed /hello endpoint
// This endpoint is always available regardless of database integration status
type HelloHandler struct {
	statusTracker StatusTrackerInterface
}

// StatusTrackerInterface defines the interface for status tracking
type StatusTrackerInterface interface {
	GetIntegrationSummary() map[string]string
	GetAllStatuses() []models.IntegrationStatus
	GetSummary() models.StatusSummary
}

// NewHelloHandler creates a new hello handler
func NewHelloHandler(tracker StatusTrackerInterface) *HelloHandler {
	return &HelloHandler{
		statusTracker: tracker,
	}
}

// Hello returns a simple greeting and server status
// @Summary Hello endpoint
// @Description Guaranteed liveness probe - always available regardless of integration status
// @Tags health
// @Produce json
// @Success 200 {object} models.HelloResponse
// @Router /hello [get]
func (h *HelloHandler) Hello(c *gin.Context) {
	response := models.HelloResponse{
		Message:   "Hello from Simple API",
		Timestamp: time.Now(),
	}

	// Add integration summary if tracker is available
	if h.statusTracker != nil {
		response.Integrations = h.statusTracker.GetIntegrationSummary()
	}

	c.JSON(http.StatusOK, response)
}