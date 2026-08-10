package main

import (
	"context"
	"log"
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

type mediaJobProcessor struct {
	mediaUseCase *usecase.MediaUseCase
}

func (m *mediaJobProcessor) ProcessJob(ctx context.Context, jobID string) error {
	// no-op for now; worker pool only needs this method to compile
	return nil
}

func main() {
	// 1. Load configuration
	cfg := config.Load()
	log.Printf("Starting server on %s", cfg.ServerAddress)

	// 2. Initialize Repositories
	jobRepository := memory.NewJobRepository()

	// 3. Initialize Frameworks (Processor)
	imageProcessor := processor.NewImageProcessor(cfg)

	// 4. Initialize Use Cases
	mediaUseCase := usecase.NewMediaUseCase(cfg, jobRepository, imageProcessor)

	// 5. Initialize Background Workers
	workerPool := worker.NewPool(4, 20, &mediaJobProcessor{mediaUseCase: mediaUseCase})
	cleanupService := worker.NewCleanupService(cfg, jobRepository)

	// --- NEW: Context for graceful shutdown ---
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start the background processes
	workerPool.Start(ctx)
	cleanupService.Start(ctx)

	// 6. Initialize Router & Server
	router := deliveryhttp.NewRouter(cfg)

	// Run the server in a goroutine so it doesn't block our shutdown listener
	go func() {
		if err := router.Run(cfg.ServerAddress); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for Ctrl+C
	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	// Stop workers
	workerPool.Stop()

	time.Sleep(1 * time.Second)
	log.Println("Goodbye!")
}
