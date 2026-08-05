package domain

import (
	"context"
	"errors"
	"time"
)

// JobStatus represents the current state of a media processing job

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

var ErrJobNotFound = errors.New("job not found")

// Job is the core entity representing an uploaded file being processed

type Job struct {
	ID           string
	OriginalName string
	OriginalSize int64
	Status       JobStatus
	Outputs      []Output
	Error        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// JobRepository defines how the application interacts with job storage.

type JobRepository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id string) (*Job, error)
	Update(ctx context.Context, job *Job) error
	Delete(ctx context.Context, id string) error
	ListExpired(ctx context.Context, before time.Time) ([]*Job, error)
}
