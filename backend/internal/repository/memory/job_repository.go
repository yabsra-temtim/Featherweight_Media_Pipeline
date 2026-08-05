package memory

import (
	"context"
	"sync"
	"time"

	"featherweight/internal/domain"
)

type JobRepository struct {
	mu   sync.RWMutex
	jobs map[string]*domain.Job
}

func NewJobRepository() *JobRepository {
	return &JobRepository{
		jobs: make(map[string]*domain.Job),
	}
}

func (r *JobRepository) Create(ctx context.Context, job *domain.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	jobCopy := *job
	r.jobs[job.ID] = &jobCopy
	return nil
}

func (r *JobRepository) GetByID(ctx context.Context, id string) (*domain.Job, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[id]
	if !exists {
		return nil, domain.ErrJobNotFound
	}
	jobCopy := *job
	return &jobCopy, nil

}

func (r *JobRepository) Update(ctx context.Context, job *domain.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[job.ID]; !exists {
		return domain.ErrJobNotFound
	}

	jobCopy := *job
	r.jobs[job.ID] = &jobCopy
	return nil
}

func (r *JobRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.jobs, id)
	return nil
}

func (r *JobRepository) ListExpired(ctx context.Context, before time.Time) ([]*domain.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var expired []*domain.Job

	for _, job := range r.jobs {
		if job.CreatedAt.Before(before) {
			jobCopy := *job
			expired = append(expired, &jobCopy)
		}
	}

	return expired, nil
}
