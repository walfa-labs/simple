package mongo

import (
	"context"
	"net/http"
	"time"

	"simple/config"
	"simple/constants"
	"simple/handlers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	collection *mongo.Collection
}

type Document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type DocumentInput struct {
	Title string `json:"title" binding:"required"`
}

// NewHandler creates a new MongoDB handler
func NewHandler(cfg config.MongoConfig) (*Handler, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBConnectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.ConnectionString()))
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := client.Database(cfg.DB).Collection("documents")

	// Migrate schema (create indexes)
	if err := MigrateSchema(collection); err != nil {
		return nil, err
	}

	return &Handler{collection: collection}, nil
}

// MigrateSchema creates required indexes
// This is called during initialization and can be called manually for upgrades
func (h *Handler) MigrateSchema() error {
	return MigrateSchema(h.collection)
}

// RegisterRoutes registers MongoDB routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/documents", h.CreateDocument)
	rg.GET("/documents", h.GetDocuments)
	rg.GET("/documents/:id", h.GetDocument)
	rg.PUT("/documents/:id", h.UpdateDocument)
	rg.DELETE("/documents/:id", h.DeleteDocument)
	rg.GET("/health", h.HealthCheck)
}

// CreateDocument creates a new document
// @Summary Create MongoDB document
// @Description Create a new document in MongoDB
// @Tags mongo
// @Accept json
// @Produce json
// @Param document body DocumentInput true "Document data"
// @Success 201 {object} Document
// @Router /mongo/documents [post]
func (h *Handler) CreateDocument(c *gin.Context) {
	var input DocumentInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc := Document{
		Title:     input.Title,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	result, err := h.collection.InsertOne(ctx, doc)
	if err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}

	doc.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, doc)
}

// GetDocuments gets all documents
// @Summary Get all MongoDB documents
// @Description Get all documents from MongoDB
// @Tags mongo
// @Produce json
// @Success 200 {array} Document
// @Router /mongo/documents [get]
func (h *Handler) GetDocuments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	cursor, err := h.collection.Find(ctx, bson.M{})
	if err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}
	defer cursor.Close(ctx)

	var documents []Document
	if err := cursor.All(ctx, &documents); err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, documents)
}

// GetDocument gets a single document
// @Summary Get MongoDB document by ID
// @Description Get a document from MongoDB by ID
// @Tags mongo
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} Document
// @Router /mongo/documents/{id} [get]
func (h *Handler) GetDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	var doc Document
	err = h.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	if err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, doc)
}

// DeleteDocument deletes a document
// @Summary Delete MongoDB document
// @Description Delete a document from MongoDB by ID
// @Tags mongo
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} map[string]string
// @Router /mongo/documents/{id} [delete]
func (h *Handler) DeleteDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	result, err := h.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "document deleted"})
}

// UpdateDocument updates a document
// @Summary Update MongoDB document
// @Description Update a document in MongoDB by ID
// @Tags mongo
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Param document body DocumentInput true "Document data"
// @Success 200 {object} Document
// @Router /mongo/documents/{id} [put]
func (h *Handler) UpdateDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input DocumentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title": input.Title,
		},
	}

	result, err := h.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	// Return updated document
	var doc Document
	if err := h.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc); err != nil {
		handlers.SanitizeDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, doc)
}

// HealthCheck checks MongoDB connection
// @Summary MongoDB health check
// @Description Check MongoDB connection status
// @Tags mongo
// @Produce json
// @Success 200 {object} map[string]string
// @Router /mongo/health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBHealthCheckTimeout)
	defer cancel()

	client := h.collection.Database().Client()
	if err := client.Ping(ctx, nil); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":           "unhealthy",
			"database":         "mongodb",
			"migration_status": "migrated",
			"error":            err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "healthy",
		"database":         "mongodb",
		"migration_status": "migrated",
	})
}
