package usecase

import (
	"context"
	"errors"
	"strings"

	"featherweight/internal/domain"
)

var ErrInvalidJobID = errors.New("job ID is required")

// JobUseCase defines the operations for the media processing jobs
type JobUseCase struct {
	repository domain.JobRepository
}

func NewJobUseCase(repository domain.JobRepository) *JobUseCase {
	return &JobUseCase{
		repository: repository,
	}
}

// CreateJob creates a new processing job and enqueues it for processing
func (u *JobUseCase) GetJob(ctx context.Context, jobID string) (*domain.Job, error) {
	jobID = strings.TrimSpace(jobID)
	if jobID == "" {
		return nil, ErrInvalidJobID
	}

	return u.repository.GetByID(ctx, jobID)

}
