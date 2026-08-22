package http

import (
	"os"
	"strings"
	"time"

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

	// ── CORS ─────────────────────────────────────────────────────────────────
	// Build the allowed-origins list from CORS_ALLOWED_ORIGINS env var
	// (comma-separated), falling back to permissive defaults for local dev.
	allowedOrigins := buildAllowedOrigins()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// ── Handlers ──────────────────────────────────────────────────────────────
	healthHandler := handler.NewHealthHandler(imageProcessor)

	uploadHandler := handler.NewUploadHandler(
		cfg,
		mediaUseCase,
		workerPool,
	)

	jobHandler := handler.NewJobHandler(
		jobUseCase,
	)

	// ── Root ──────────────────────────────────────────────────────────────────
	// Returns 200 so Render's health check on / passes and visitors get
	// a clear message instead of a blank 404.
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "Featherweight Media Pipeline API",
			"version": "v1",
			"status":  "ok",
			"docs":    "/api/v1/health",
		})
	})

	// ── API routes ────────────────────────────────────────────────────────────
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)
		v1.POST("/media/upload", uploadHandler.Upload)
		v1.GET("/jobs/:id", jobHandler.GetByID)
	}

	return router
}

// buildAllowedOrigins returns the CORS origin whitelist.
// Set CORS_ALLOWED_ORIGINS to a comma-separated list in your Render dashboard,
// e.g. "https://featherweight-ui.onrender.com"
func buildAllowedOrigins() []string {
	if raw := os.Getenv("CORS_ALLOWED_ORIGINS"); raw != "" {
		var origins []string
		for _, o := range strings.Split(raw, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				origins = append(origins, trimmed)
			}
		}
		if len(origins) > 0 {
			return origins
		}
	}
	// Local dev fallback — allow all localhost variants
	return []string{
		"http://localhost:5173",
		"http://localhost:4173",
		"http://localhost:3000",
	}
}
