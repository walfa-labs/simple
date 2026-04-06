package postgres

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"simple/config"
	"simple/constants"
	"simple/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Handler struct {
	db *sql.DB
}

type Record struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type RecordInput struct {
	Title string `json:"title" binding:"required"`
}

// NewHandler creates a new PostgreSQL handler
func NewHandler(cfg config.PostgresConfig) (*Handler, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, err
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBConnectionTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	// Migrate schema (create tables and indexes)
	if err := MigrateSchema(db); err != nil {
		return nil, err
	}

	return &Handler{db: db}, nil
}

// MigrateSchema creates required database schema
// This is called during initialization and can be called manually for upgrades
func (h *Handler) MigrateSchema() error {
	return MigrateSchema(h.db)
}

// RegisterRoutes registers PostgreSQL routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/records", h.CreateRecord)
	rg.GET("/records", h.GetRecords)
	rg.GET("/records/:id", h.GetRecord)
	rg.PUT("/records/:id", h.UpdateRecord)
	rg.DELETE("/records/:id", h.DeleteRecord)
	rg.GET("/health", h.HealthCheck)
}

// CreateRecord creates a new record
// @Summary Create PostgreSQL record
// @Description Create a new record in PostgreSQL
// @Tags postgres
// @Accept json
// @Produce json
// @Param record body RecordInput true "Record data"
// @Success 201 {object} Record
// @Router /postgres/records [post]
func (h *Handler) CreateRecord(c *gin.Context) {
	var input RecordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var record Record
	err := h.db.QueryRow(
		"INSERT INTO records (title) VALUES ($1) RETURNING id, title, created_at",
		input.Title,
	).Scan(&record.ID, &record.Title, &record.CreatedAt)

	if err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}

	c.JSON(http.StatusCreated, record)
}

// GetRecords gets all records
// @Summary Get all PostgreSQL records
// @Description Get all records from PostgreSQL
// @Tags postgres
// @Produce json
// @Success 200 {array} Record
// @Router /postgres/records [get]
func (h *Handler) GetRecords(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	rows, err := h.db.QueryContext(ctx, "SELECT id, title, created_at FROM records ORDER BY created_at DESC")
	if err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.ID, &record.Title, &record.CreatedAt); err != nil {
			handlers.SanitizeDBError(c, err)
			return
		}
		records = append(records, record)
	}

	c.JSON(http.StatusOK, records)
}

// GetRecord gets a single record
// @Summary Get PostgreSQL record by ID
// @Description Get a record from PostgreSQL by ID
// @Tags postgres
// @Produce json
// @Param id path int true "Record ID"
// @Success 200 {object} Record
// @Router /postgres/records/{id} [get]
func (h *Handler) GetRecord(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	var record Record
	err := h.db.QueryRowContext(ctx,
		"SELECT id, title, created_at FROM records WHERE id = $1",
		id,
	).Scan(&record.ID, &record.Title, &record.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	if err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteRecord deletes a record
// @Summary Delete PostgreSQL record
// @Description Delete a record from PostgreSQL by ID
// @Tags postgres
// @Produce json
// @Param id path int true "Record ID"
// @Success 200 {object} map[string]string
// @Router /postgres/records/{id} [delete]
func (h *Handler) DeleteRecord(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	result, err := h.db.ExecContext(ctx, "DELETE FROM records WHERE id = $1", id)
	if err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "record deleted"})
}

// UpdateRecord updates a record
// @Summary Update PostgreSQL record
// @Description Update a record in PostgreSQL by ID
// @Tags postgres
// @Accept json
// @Produce json
// @Param id path int true "Record ID"
// @Param record body RecordInput true "Record data"
// @Success 200 {object} Record
// @Router /postgres/records/{id} [put]
func (h *Handler) UpdateRecord(c *gin.Context) {
	id := c.Param("id")

	var input RecordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	var record Record
	err := h.db.QueryRowContext(ctx,
		"UPDATE records SET title = $1 WHERE id = $2 RETURNING id, title, created_at",
		input.Title, id,
	).Scan(&record.ID, &record.Title, &record.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	if err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, record)
}

// HealthCheck checks PostgreSQL connection
// @Summary PostgreSQL health check
// @Description Check PostgreSQL connection status
// @Tags postgres
// @Produce json
// @Success 200 {object} map[string]string
// @Router /postgres/health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBHealthCheckTimeout)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":           "unhealthy",
			"database":         "postgresql",
			"migration_status": "migrated",
			"error":            err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "healthy",
		"database":         "postgresql",
		"migration_status": "migrated",
	})
}
