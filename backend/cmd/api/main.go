package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"featherweight/internal/config"
	deliveryhttp "featherweight/internal/delivery/http"
	"featherweight/internal/domain"
	infracloudinary "featherweight/internal/infrastructure/cloudinary"
	infrapostgres "featherweight/internal/infrastructure/postgres"
	repopostgres "featherweight/internal/infrastructure/postgres"
	"featherweight/internal/processor"
	"featherweight/internal/usecase"
	"featherweight/internal/worker"
)

func main() {
	// 0. Load .env file
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("No .env file found, reading environment variables directly")
		}
	}

	// 1. Load configuration
	cfg := config.Load()

	log.Printf(
		"starting server on %s (store: %s)",
		cfg.ServerAddress,
		cfg.JobStore,
	)

	// 2. Create shutdown context
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	// 3. Initialize PostgreSQL
	var dbPool *pgxpool.Pool

	var err error
	dbPool, err = infrapostgres.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}
	defer dbPool.Close()

	log.Println("PostgreSQL connection established")

	// 4. Run database migrations
	if err := infrapostgres.Migrate(ctx, dbPool); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("database migrations completed")

	// 5. Initialize PostgreSQL repository
	var jobRepository domain.JobRepository

	jobRepository = repopostgres.NewJobRepository(dbPool)

	log.Printf(
		"job repository initialized: %T",
		jobRepository,
	)

	// 6. Initialize image processor
	imageProcessor := processor.NewImageProcessor(cfg)

	// 7. Initialize Cloudinary uploader
	cloudUploader, err := infracloudinary.NewUploader(cfg)
	if err != nil {
		log.Fatalf("failed to init Cloudinary: %v", err)
	}

	log.Println("Cloudinary uploader initialized")

	// 8. Initialize use cases
	jobUseCase := usecase.NewJobUseCase(
		jobRepository,
	)

	mediaUseCase := usecase.NewMediaUseCase(
		cfg,
		jobRepository,
		imageProcessor,
		cloudUploader,
	)

	// 9. Initialize background workers
	workerPool := worker.NewPool(
		cfg.WorkerCount,
		20,
		mediaUseCase,
	)

	cleanupService := worker.NewCleanupService(
		cfg,
		jobRepository,
		cloudUploader,
	)

	// 10. Start background workers
	workerPool.Start(ctx)
	cleanupService.Start(ctx)

	// 11. Initialize HTTP router
	router := deliveryhttp.NewRouter(
		cfg,
		imageProcessor,
		mediaUseCase,
		jobUseCase,
		workerPool,
	)

	// 12. Create HTTP server
	server := &http.Server{
		Addr:              cfg.ServerAddress,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      5 * time.Minute,
		IdleTimeout:       60 * time.Second,
	}

	// 13. Start HTTP server
	serverError := make(chan error, 1)

	go func() {
		log.Printf("listening on %s", cfg.ServerAddress)

		if err := server.ListenAndServe(); err != nil {
			serverError <- err
		}
	}()

	// 14. Wait for server failure or shutdown signal
	select {
	case err := <-serverError:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}

	case <-ctx.Done():
		log.Println("shutdown signal received")
	}

	// 15. Gracefully shut down HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	// 16. Stop background workers
	workerPool.Stop()

	log.Println("server stopped")
}
