package constants

import "time"

// Database timeouts
const (
	// DBConnectionTimeout is the timeout for establishing database connections
	DBConnectionTimeout = 10 * time.Second

	// DBQueryTimeout is the timeout for database queries
	DBQueryTimeout = 5 * time.Second

	// DBHealthCheckTimeout is the timeout for health check pings
	DBHealthCheckTimeout = 2 * time.Second
)

// HTTP timeouts
const (
	// HTTPReadTimeout is the timeout for reading HTTP requests
	HTTPReadTimeout = 5 * time.Second

	// HTTPWriteTimeout is the timeout for writing HTTP responses
	HTTPWriteTimeout = 5 * time.Second
)
