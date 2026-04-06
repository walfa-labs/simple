package models

import "time"

// MigrationStatus represents the state of schema migration
type MigrationStatus string

const (
	MigrationStatusMigrated    MigrationStatus = "migrated"
	MigrationStatusFailed      MigrationStatus = "failed"
	MigrationStatusNotConfigured MigrationStatus = "not_configured"
)

// IntegrationStatus tracks the operational state of each database integration
type IntegrationStatus struct {
	Name            string          // Integration identifier ("postgres", "redis", "mongo")
	Enabled         bool            // Whether integration is operational
	Error           string          // Failure reason if disabled
	LastCheckTime   time.Time       // When status was last updated
	MigrationStatus MigrationStatus // Schema migration status
}

// IntegrationSummary provides a simplified status for hello response
type IntegrationSummary struct {
	Name    string // Integration identifier
	Status  string // "enabled", "disabled", or "not_configured"
}

// HealthCheckResponse is the response format for health check endpoints
type HealthCheckResponse struct {
	Status          string          `json:"status"`
	Database        string          `json:"database"`
	MigrationStatus MigrationStatus `json:"migration_status"`
	Error           string          `json:"error,omitempty"`
}

// HelloResponse is the response for the guaranteed /hello endpoint
type HelloResponse struct {
	Message     string                       `json:"message"`
	Timestamp   time.Time                    `json:"timestamp"`
	Integrations map[string]string           `json:"integrations,omitempty"`
}

// IntegrationsStatusResponse is the response for /integrations/status endpoint
type IntegrationsStatusResponse struct {
	Integrations []IntegrationStatus `json:"integrations"`
	Summary      StatusSummary       `json:"summary"`
}

// StatusSummary provides aggregate counts
type StatusSummary struct {
	Total     int `json:"total"`
	Enabled   int `json:"enabled"`
	Disabled  int `json:"disabled"`
}

// IntegrationNames are the valid integration identifiers
const (
	IntegrationPostgres = "postgres"
	IntegrationRedis    = "redis"
	IntegrationMongo    = "mongo"
)

// ValidIntegrationNames contains all valid integration identifiers
var ValidIntegrationNames = []string{
	IntegrationPostgres,
	IntegrationRedis,
	IntegrationMongo,
}

// IsValidIntegrationName checks if the name is valid
func IsValidIntegrationName(name string) bool {
	for _, valid := range ValidIntegrationNames {
		if valid == name {
			return true
		}
	}
	return false
}