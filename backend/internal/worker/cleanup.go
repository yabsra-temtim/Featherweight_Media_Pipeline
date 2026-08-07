package worker

import (
	"context"
	"log"
	"os"
	"time"

	"featherweight/internal/config"
	"featherweight/internal/domain"
)

type CleanupService struct {
	config     config.Config
	repository domain.JobRepository
}

func NewCleanupService(cfg config.Config, repo domain.JobRepository) *CleanupService {
	return &CleanupService{
		config:     cfg,
		repository: repo,
	}
}

func (s *CleanupService) Start(ctx context.Context) {
	go func() {
		// Run every 10 minutes by default
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Cleanup service stopped")
				return
			case <-ticker.C:
				s.cleanup(ctx)
			}
		}
	}()
}

func (s *CleanupService) cleanup(ctx context.Context) {
	// Jobs older than 1 hour will be deleted
	before := time.Now().Add(-1 * time.Hour)

	expiredJobs, err := s.repository.ListExpired(ctx, before)
	if err != nil {
		log.Printf("[Cleanup] Failed to list expired jobs: %v", err)
		return
	}

	for _, job := range expiredJobs {
		// 1. Delete processed outputs
		for _, output := range job.Outputs {
			if err := os.Remove(output.Path); err != nil {
				log.Printf("[Cleanup] Failed to delete output %s: %v", output.Path, err)
			}
		}

		// 3. Remove from database
		if err := s.repository.Delete(ctx, job.ID); err != nil {
			log.Printf("[Cleanup] Failed to delete job %s from repo: %v", job.ID, err)
		} else {
			log.Printf("[Cleanup] Successfully cleaned up job %s", job.ID)
		}
	}
}
