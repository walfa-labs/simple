package redis

// MigrateSchema validates Redis connection (Redis has no schema)
// This is a no-op function for consistency across integrations
func MigrateSchema() error {
	// Redis is a key-value store with no traditional schema
	// Connection validation is done in NewHandler via Ping
	// No additional migration needed
	return nil
}

// GetSchemaVersion returns the current schema version (empty for Redis)
func GetSchemaVersion() string {
	return "" // Redis has no schema version
}