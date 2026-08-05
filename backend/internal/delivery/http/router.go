package http

import (
	"featherweight/internal/config"
	"featherweight/internal/delivery/handler"

	"github.com/gin-gonic/gin"
)

// NewRouter initializes the Gin engine and registers all routes
func NewRouter(cfg config.Config) *gin.Engine {
	// Initialize Gin with default Logger and Recovery middleware
	router := gin.Default()

	// Initialize our handlers
	healthHandler := handler.NewHealthHandler()

	// Set up the API routing group
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)
	}

	return router
}
