package handlers

import (
	"net/http"

	"simple/models"

	"github.com/gin-gonic/gin"
)

// IntegrationHandler handles integration status endpoints
type IntegrationHandler struct {
	statusTracker StatusTrackerInterface
}

// NewIntegrationHandler creates a new integration handler
func NewIntegrationHandler(tracker StatusTrackerInterface) *IntegrationHandler {
	return &IntegrationHandler{
		statusTracker: tracker,
	}
}

// GetStatus returns the status of all integrations
// @Summary Integration status
// @Description Get operational status of all database integrations
// @Tags integrations
// @Produce json
// @Success 200 {object} models.IntegrationsStatusResponse
// @Router /integrations/status [get]
func (h *IntegrationHandler) GetStatus(c *gin.Context) {
	response := models.IntegrationsStatusResponse{
		Integrations: h.statusTracker.GetAllStatuses(),
		Summary:      h.statusTracker.GetSummary(),
	}

	c.JSON(http.StatusOK, response)
}

// DisabledIntegrationResponse returns error for disabled integration endpoints
func DisabledIntegrationResponse(c *gin.Context, integrationName string) {
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error":            "integration disabled",
		"integration":      integrationName,
		"migration_status": "failed or not_configured",
	})
}