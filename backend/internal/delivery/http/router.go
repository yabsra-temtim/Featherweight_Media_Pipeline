package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"featherweight/internal/config"
	"featherweight/internal/delivery/handler"
	"featherweight/internal/processor"
	"featherweight/internal/usecase"
	"featherweight/internal/worker"
)

func NewRouter(
	cfg config.Config,
	imageProcessor *processor.ImageProcessor,
	mediaUseCase *usecase.MediaUseCase,
	jobUseCase *usecase.JobUseCase,
	workerPool *worker.Pool,
) *gin.Engine {

	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.Default())

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(imageProcessor)

	uploadHandler := handler.NewUploadHandler(
		cfg,
		mediaUseCase,
		workerPool,
	)

	jobHandler := handler.NewJobHandler(
		jobUseCase,
	)

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)
		v1.POST("/media/upload", uploadHandler.Upload)
		v1.GET("/jobs/:id", jobHandler.GetByID)
	}

	return router
}
