package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"featherweight/internal/domain"
)

type JobRepository struct {
	pool *pgxpool.Pool
}

func NewJobRepository(pool *pgxpool.Pool) *JobRepository {
	return &JobRepository{
		pool: pool,
	}
}

func (r *JobRepository) Create(ctx context.Context, job *domain.Job) error {
	formatsJSON, _ := json.Marshal(job.Formats)
	outputsJSON, _ := json.Marshal(job.Outputs)

	query := `
    INSERT INTO jobs (id, original_name, original_path, original_size, width, quality, formats, status, outputs, error, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
  `
	_, err := r.pool.Exec(ctx, query,
		job.ID, job.OriginalName, job.OriginalPath, job.OriginalSize, job.Width, job.Quality,
		formatsJSON, job.Status, outputsJSON, job.Error, job.CreatedAt, job.UpdatedAt)
	return err
}

func (r *JobRepository) GetByID(ctx context.Context, id string) (*domain.Job, error) {
	query := `
    SELECT id, original_name, original_path, original_size, width, quality, formats, status, outputs, error, created_at, updated_at
    FROM jobs
    WHERE id = $1
  `
	row := r.pool.QueryRow(ctx, query, id)

	var job domain.Job
	var formatsJSON, outputsJSON []byte

	err := row.Scan(
		&job.ID, &job.OriginalName, &job.OriginalPath, &job.OriginalSize, &job.Width, &job.Quality,
		&formatsJSON, &job.Status, &outputsJSON, &job.Error, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrJobNotFound
		}
		return nil, err
	}

	json.Unmarshal(formatsJSON, &job.Formats)
	json.Unmarshal(outputsJSON, &job.Outputs)

	return &job, nil
}

func (r *JobRepository) Update(ctx context.Context, job *domain.Job) error {
	outputsJSON, _ := json.Marshal(job.Outputs)

	query := `
    UPDATE jobs
    SET status = $1, outputs = $2, error = $3, updated_at = $4
    WHERE id = $5
  `
	_, err := r.pool.Exec(ctx, query, job.Status, outputsJSON, job.Error, job.UpdatedAt, job.ID)
	return err
}

func (r *JobRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM jobs WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *JobRepository) ListExpired(ctx context.Context, before time.Time) ([]*domain.Job, error) {
	query := `
    SELECT id, original_name, original_path, original_size, width, quality, formats, status, outputs, error, created_at, updated_at
    FROM jobs
    WHERE created_at < $1
  `
	rows, err := r.pool.Query(ctx, query, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expired []*domain.Job
	for rows.Next() {
		var job domain.Job
		var formatsJSON, outputsJSON []byte

		if err := rows.Scan(
			&job.ID, &job.OriginalName, &job.OriginalPath, &job.OriginalSize, &job.Width, &job.Quality,
			&formatsJSON, &job.Status, &outputsJSON, &job.Error, &job.CreatedAt, &job.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(formatsJSON, &job.Formats)
		json.Unmarshal(outputsJSON, &job.Outputs)
		expired = append(expired, &job)
	}
	return expired, rows.Err()
}
