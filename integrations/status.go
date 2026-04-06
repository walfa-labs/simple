package integrations

import (
	"sync"
	"time"

	"simple/models"
)

// StatusTracker tracks the operational status of all integrations
type StatusTracker struct {
	mu     sync.RWMutex
	statuses map[string]*models.IntegrationStatus
}

// NewStatusTracker creates a new StatusTracker
func NewStatusTracker() *StatusTracker {
	tracker := &StatusTracker{
		statuses: make(map[string]*models.IntegrationStatus),
	}
	// Initialize all integrations with not_configured status
	for _, name := range models.ValidIntegrationNames {
		tracker.statuses[name] = &models.IntegrationStatus{
			Name:            name,
			Enabled:         false,
			Error:           "",
			LastCheckTime:   time.Now(),
			MigrationStatus: models.MigrationStatusNotConfigured,
		}
	}
	return tracker
}

// SetStatus updates the status for an integration
func (t *StatusTracker) SetStatus(name string, status *models.IntegrationStatus) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !models.IsValidIntegrationName(name) {
		return
	}

	status.LastCheckTime = time.Now()
	t.statuses[name] = status
}

// SetSuccess marks an integration as successfully migrated
func (t *StatusTracker) SetSuccess(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !models.IsValidIntegrationName(name) {
		return
	}

	t.statuses[name] = &models.IntegrationStatus{
		Name:            name,
		Enabled:         true,
		Error:           "",
		LastCheckTime:   time.Now(),
		MigrationStatus: models.MigrationStatusMigrated,
	}
}

// SetFailure marks an integration as failed
func (t *StatusTracker) SetFailure(name string, errorMsg string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !models.IsValidIntegrationName(name) {
		return
	}

	// Truncate error if too long
	if len(errorMsg) > 500 {
		errorMsg = errorMsg[:500]
	}

	t.statuses[name] = &models.IntegrationStatus{
		Name:            name,
		Enabled:         false,
		Error:           errorMsg,
		LastCheckTime:   time.Now(),
		MigrationStatus: models.MigrationStatusFailed,
	}
}

// SetNotConfigured marks an integration as not configured
func (t *StatusTracker) SetNotConfigured(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !models.IsValidIntegrationName(name) {
		return
	}

	t.statuses[name] = &models.IntegrationStatus{
		Name:            name,
		Enabled:         false,
		Error:           "",
		LastCheckTime:   time.Now(),
		MigrationStatus: models.MigrationStatusNotConfigured,
	}
}

// GetStatus returns the status for a specific integration
func (t *StatusTracker) GetStatus(name string) *models.IntegrationStatus {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if status, exists := t.statuses[name]; exists {
		// Return a copy to prevent modification
		copy := *status
		return &copy
	}
	return nil
}

// GetAllStatuses returns all integration statuses
func (t *StatusTracker) GetAllStatuses() []models.IntegrationStatus {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]models.IntegrationStatus, 0, len(t.statuses))
	for _, name := range models.ValidIntegrationNames {
		if status, exists := t.statuses[name]; exists {
			result = append(result, *status)
		}
	}
	return result
}

// GetSummary returns aggregate status counts
func (t *StatusTracker) GetSummary() models.StatusSummary {
	t.mu.RLock()
	defer t.mu.RUnlock()

	summary := models.StatusSummary{
		Total: len(t.statuses),
	}

	for _, status := range t.statuses {
		if status.Enabled {
			summary.Enabled++
		} else {
			summary.Disabled++
		}
	}

	return summary
}

// GetIntegrationSummary returns a map of integration names to their status strings
func (t *StatusTracker) GetIntegrationSummary() map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make(map[string]string)
	for _, name := range models.ValidIntegrationNames {
		if status, exists := t.statuses[name]; exists {
			if status.Enabled {
				result[name] = "enabled"
			} else if status.MigrationStatus == models.MigrationStatusNotConfigured {
				result[name] = "not_configured"
			} else {
				result[name] = "disabled"
			}
		}
	}
	return result
}

// IsEnabled checks if a specific integration is enabled
func (t *StatusTracker) IsEnabled(name string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if status, exists := t.statuses[name]; exists {
		return status.Enabled
	}
	return false
}