package handlers

import (
	"net/http"

	"simple/models"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles global health check
type HealthHandler struct {
	statusTracker StatusTrackerInterface
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(tracker StatusTrackerInterface) *HealthHandler {
	return &HealthHandler{
		statusTracker: tracker,
	}
}

// HealthCheck returns the overall health status
// @Summary Health check
// @Description Check overall API health with integration statuses
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	integrationStatuses := h.statusTracker.GetAllStatuses()

	// Build integration map
	integrations := make(map[string]interface{})
	overallHealthy := false

	for _, status := range integrationStatuses {
		integrationData := map[string]interface{}{
			"status":           "disabled",
			"migration_status": status.MigrationStatus,
		}

		if status.Enabled {
			integrationData["status"] = "healthy"
			overallHealthy = true
		} else if status.MigrationStatus == models.MigrationStatusFailed {
			integrationData["status"] = "unhealthy"
			integrationData["error"] = status.Error
		} else if status.MigrationStatus == models.MigrationStatusNotConfigured {
			integrationData["status"] = "not_configured"
		}

		integrations[status.Name] = integrationData
	}

	// Overall status is healthy if at least one integration is working
	// or if the server is running (even with no integrations)
	overallStatus := "healthy"
	if !overallHealthy && len(integrationStatuses) > 0 {
		// Check if all are not_configured (server still healthy)
		allNotConfigured := true
		for _, status := range integrationStatuses {
			if status.MigrationStatus != models.MigrationStatusNotConfigured {
				allNotConfigured = false
				break
			}
		}
		if !allNotConfigured {
			overallStatus = "degraded"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       overallStatus,
		"integrations": integrations,
	})
}
