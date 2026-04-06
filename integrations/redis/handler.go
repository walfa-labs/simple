package redis

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"simple/config"
	"simple/constants"
	"simple/handlers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	client *redis.Client
}

type CacheItem struct {
	Key       string      `json:"key"`
	Value     string      `json:"value"`
	ExpiresIn int         `json:"expires_in,omitempty"` // seconds
	Data      interface{} `json:"data,omitempty"`
}

// NewHandler creates a new Redis handler
func NewHandler(cfg config.RedisConfig) (*Handler, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: cfg.Password,
		DB:       0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBConnectionTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	// Migrate schema (no-op for Redis, connection validated)
	if err := MigrateSchema(); err != nil {
		return nil, err
	}

	return &Handler{client: client}, nil
}

// MigrateSchema validates schema (no-op for Redis)
// This is included for consistency with other integrations
func (h *Handler) MigrateSchema() error {
	return MigrateSchema()
}

// RegisterRoutes registers Redis routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/cache", h.SetCache)
	rg.GET("/cache/:key", h.GetCache)
	rg.DELETE("/cache/:key", h.DeleteCache)
	rg.GET("/health", h.HealthCheck)
}

// SetCache sets a cache item
// @Summary Set Redis cache item
// @Description Set a key-value pair in Redis cache
// @Tags redis
// @Accept json
// @Produce json
// @Param item body CacheItem true "Cache item data"
// @Success 200 {object} map[string]string
// @Router /redis/cache [post]
func (h *Handler) SetCache(c *gin.Context) {
	var item CacheItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	var err error
	if item.ExpiresIn > 0 {
		err = h.client.Set(ctx, item.Key, item.Value, time.Duration(item.ExpiresIn)*time.Second).Err()
	} else {
		err = h.client.Set(ctx, item.Key, item.Value, 0).Err()
	}

	if err != nil {
		handlers.SanitizeCacheError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cache set", "key": item.Key})
}

// GetCache gets a cache item
// @Summary Get Redis cache item
// @Description Get a value from Redis cache by key
// @Tags redis
// @Produce json
// @Param key path string true "Cache key"
// @Success 200 {object} CacheItem
// @Router /redis/cache/{key} [get]
func (h *Handler) GetCache(c *gin.Context) {
	key := c.Param("key")

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	val, err := h.client.Get(ctx, key).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}
	if err != nil {
		handlers.SanitizeCacheError(c, err)
		return
	}

	// Get TTL
	ttl, err := h.client.TTL(ctx, key).Result()
	if err != nil {
		handlers.SanitizeCacheError(c, err)
		return
	}

	item := CacheItem{
		Key:   key,
		Value: val,
	}

	if ttl > 0 {
		item.ExpiresIn = int(ttl.Seconds())
	}

	c.JSON(http.StatusOK, item)
}

// DeleteCache deletes a cache item
// @Summary Delete Redis cache item
// @Description Delete a key from Redis cache
// @Tags redis
// @Produce json
// @Param key path string true "Cache key"
// @Success 200 {object} map[string]string
// @Router /redis/cache/{key} [delete]
func (h *Handler) DeleteCache(c *gin.Context) {
	key := c.Param("key")

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	if err := h.client.Del(ctx, key).Err(); err != nil {
		handlers.SanitizeCacheError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cache deleted", "key": key})
}

// HealthCheck checks Redis connection
// @Summary Redis health check
// @Description Check Redis connection status
// @Tags redis
// @Produce json
// @Success 200 {object} map[string]string
// @Router /redis/health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBHealthCheckTimeout)
	defer cancel()

	if err := h.client.Ping(ctx).Err(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":           "unhealthy",
			"database":         "redis",
			"migration_status": "migrated",
			"error":            err.Error(),
		})
		return
	}

	info, err := h.client.Info(ctx, "server").Result()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":           "unhealthy",
			"database":         "redis",
			"migration_status": "migrated",
			"error":            err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "healthy",
		"database":         "redis",
		"migration_status": "migrated",
		"info":             strconv.Itoa(len(info)) + " bytes",
	})
}
