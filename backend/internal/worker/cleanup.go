package worker

import (
	"context"
	"log"
	"time"

	"featherweight/internal/config"
	"featherweight/internal/domain"
)

// CloudinaryDeleter is the minimal interface the cleanup service needs.
type CloudinaryDeleter interface {
	Delete(ctx context.Context, publicID string) error
}

type CleanupService struct {
	config     config.Config
	repository domain.JobRepository
	cloud      CloudinaryDeleter
}

func NewCleanupService(cfg config.Config, repo domain.JobRepository, cloud CloudinaryDeleter) *CleanupService {
	return &CleanupService{
		config:     cfg,
		repository: repo,
		cloud:      cloud,
	}
}

func (s *CleanupService) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(s.config.CleanupInterval)
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
	before := time.Now().Add(-s.config.FileRetention)

	expiredJobs, err := s.repository.ListExpired(ctx, before)
	if err != nil {
		log.Printf("[Cleanup] Failed to list expired jobs: %v", err)
		return
	}

	for _, job := range expiredJobs {
		// 1. Delete each output from Cloudinary using its public_id.
		for _, output := range job.Outputs {
			if output.PublicID == "" {
				continue
			}
			if err := s.cloud.Delete(ctx, output.PublicID); err != nil {
				log.Printf("[Cleanup] Failed to delete Cloudinary asset %s: %v", output.PublicID, err)
			}
		}

		// 2. Remove from database.
		if err := s.repository.Delete(ctx, job.ID); err != nil {
			log.Printf("[Cleanup] Failed to delete job %s from repo: %v", job.ID, err)
		} else {
			log.Printf("[Cleanup] Successfully cleaned up job %s", job.ID)
		}
	}
}
