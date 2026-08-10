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

	"featherweight/internal/config"
	deliveryhttp "featherweight/internal/delivery/http"
	"featherweight/internal/processor"
	"featherweight/internal/repository/memory"
	"featherweight/internal/usecase"
	"featherweight/internal/worker"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()

	log.Printf(
		"configuration loaded: server=%s",
		cfg.ServerAddress,
	)

	// 2. Initialize repository
	jobRepository := memory.NewJobRepository()

	log.Printf(
		"job repository initialized: %T",
		jobRepository,
	)

	// 3. Initialize image processor
	imageProcessor := processor.NewImageProcessor(cfg)

	// 4. Initialize use cases
	jobUseCase := usecase.NewJobUseCase(
		jobRepository,
	)

	mediaUseCase := usecase.NewMediaUseCase(
		cfg,
		jobRepository,
		imageProcessor,
	)

	// 5. Initialize background workers
	workerPool := worker.NewPool(
		4,
		20,
		mediaUseCase,
	)

	cleanupService := worker.NewCleanupService(
		cfg,
		jobRepository,
	)

	// 6. Create shutdown context
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	// 7. Start background workers
	workerPool.Start(ctx)
	cleanupService.Start(ctx)

	// 8. Initialize Gin router
	router := deliveryhttp.NewRouter(
		cfg,
		imageProcessor,
		mediaUseCase,
		jobUseCase,
		workerPool,
	)

	// 9. Create HTTP server
	server := &http.Server{
		Addr:              cfg.ServerAddress,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      5 * time.Minute,
		IdleTimeout:       60 * time.Second,
	}

	// 10. Start HTTP server
	serverError := make(chan error, 1)

	go func() {
		log.Printf(
			"listening on %s",
			cfg.ServerAddress,
		)

		serverError <- server.ListenAndServe()
	}()

	// 11. Wait for server failure or shutdown signal
	select {
	case err := <-serverError:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(
				"server failed: %v",
				err,
			)
		}

	case <-ctx.Done():
		log.Println("shutdown signal received")
	}

	// 12. Gracefully shut down HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf(
			"server shutdown failed: %v",
			err,
		)
	}

	// 13. Stop workers
	workerPool.Stop()

	log.Println("server stopped")
}
