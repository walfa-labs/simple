package constants

// IntegrationNames are the valid database integration identifiers
const (
	IntegrationPostgres = "postgres"
	IntegrationRedis    = "redis"
	IntegrationMongo    = "mongo"
)

// AllIntegrationNames contains all valid integration identifiers
var AllIntegrationNames = []string{
	IntegrationPostgres,
	IntegrationRedis,
	IntegrationMongo,
}

// IsValidIntegrationName checks if the name is a valid integration
func IsValidIntegrationName(name string) bool {
	for _, valid := range AllIntegrationNames {
		if valid == name {
			return true
		}
	}
	return false
}