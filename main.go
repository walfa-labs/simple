package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"simple/config"
	"simple/handlers"
	"simple/integrations"
	"simple/integrations/mongo"
	"simple/integrations/postgres"
	"simple/integrations/redis"
	"simple/models"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "simple/docs"
)

// @title Simple Go API
// @version 1.0
// @description A simple Go API with optional PostgreSQL, Redis, and MongoDB integrations
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize status tracker for all integrations
	statusTracker := integrations.NewStatusTracker()

	// Initialize guaranteed endpoints (always available)
	helloHandler := handlers.NewHelloHandler(statusTracker)
	router.GET("/hello", helloHandler.Hello)

	integrationHandler := handlers.NewIntegrationHandler(statusTracker)
	router.GET("/integrations/status", integrationHandler.GetStatus)

	// Initialize integrations
	var (
		pgHandler    *postgres.Handler
		redisHandler *redis.Handler
		mongoHandler *mongo.Handler
	)

	// PostgreSQL initialization with graceful degradation
	if cfg.Postgres.Enabled {
		pgHandler, err = postgres.NewHandler(cfg.Postgres)
		if err != nil {
			log.Printf("[POSTGRES] migration failed: %v", err)
			statusTracker.SetFailure(models.IntegrationPostgres, err.Error())
			// Register disabled routes
			registerDisabledRoutes(router, models.IntegrationPostgres)
		} else {
			log.Println("PostgreSQL integration enabled")
			statusTracker.SetSuccess(models.IntegrationPostgres)
			pgHandler.RegisterRoutes(router.Group("/postgres"))
		}
	} else {
		log.Println("PostgreSQL integration disabled (missing .env config)")
		statusTracker.SetNotConfigured(models.IntegrationPostgres)
	}

	// Redis initialization with graceful degradation
	if cfg.Redis.Enabled {
		redisHandler, err = redis.NewHandler(cfg.Redis)
		if err != nil {
			log.Printf("[REDIS] migration failed: %v", err)
			statusTracker.SetFailure(models.IntegrationRedis, err.Error())
			// Register disabled routes
			registerDisabledRoutes(router, models.IntegrationRedis)
		} else {
			log.Println("Redis integration enabled")
			statusTracker.SetSuccess(models.IntegrationRedis)
			redisHandler.RegisterRoutes(router.Group("/redis"))
		}
	} else {
		log.Println("Redis integration disabled (missing .env config)")
		statusTracker.SetNotConfigured(models.IntegrationRedis)
	}

	// MongoDB initialization with graceful degradation
	if cfg.Mongo.Enabled {
		mongoHandler, err = mongo.NewHandler(cfg.Mongo)
		if err != nil {
			log.Printf("[MONGO] migration failed: %v", err)
			statusTracker.SetFailure(models.IntegrationMongo, err.Error())
			// Register disabled routes
			registerDisabledRoutes(router, models.IntegrationMongo)
		} else {
			log.Println("MongoDB integration enabled")
			statusTracker.SetSuccess(models.IntegrationMongo)
			mongoHandler.RegisterRoutes(router.Group("/mongo"))
		}
	} else {
		log.Println("MongoDB integration disabled (missing .env config)")
		statusTracker.SetNotConfigured(models.IntegrationMongo)
	}

	// Health check endpoint with integration status
	healthHandler := handlers.NewHealthHandler(statusTracker)
	router.GET("/health", healthHandler.HealthCheck)

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// registerDisabledRoutes registers catch-all routes that return 503 for disabled integrations
func registerDisabledRoutes(router *gin.Engine, integrationName string) {
	router.Any("/"+integrationName+"/*action", func(c *gin.Context) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":            "integration disabled",
			"integration":      integrationName,
			"migration_status": "failed",
		})
	})
}