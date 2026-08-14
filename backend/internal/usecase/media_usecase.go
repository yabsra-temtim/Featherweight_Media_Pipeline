package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"featherweight/internal/config"
	"featherweight/internal/domain"
	"featherweight/internal/processor"

	"github.com/google/uuid"
)

var (
	ErrNoFile       = errors.New("an image file is required")
	ErrInvalidWidth = errors.New("width must be between 1 and the configured maximum")
)

// CloudinaryStore is the interface for cloud storage operations.
// Implemented by infrastructure/cloudinary.Uploader.
type CloudinaryStore interface {
	Upload(ctx context.Context, localPath, publicID string) (secureURL string, fullPublicID string, err error)
	Delete(ctx context.Context, publicID string) error
}

// UploadInput represents the data expected from the HTTP handler.
type UploadInput struct {
	FileName string
	FileSize int64
	File     io.Reader
	Width    int
	Quality  int
	Formats  []string
}

// MediaUseCase orchestrates upload and image processing.
type MediaUseCase struct {
	config     config.Config
	repository domain.JobRepository
	processor  *processor.ImageProcessor
	cloud      CloudinaryStore
}

func NewMediaUseCase(
	cfg config.Config,
	repo domain.JobRepository,
	proc *processor.ImageProcessor,
	cloud CloudinaryStore,
) *MediaUseCase {
	return &MediaUseCase{
		config:     cfg,
		repository: repo,
		processor:  proc,
		cloud:      cloud,
	}
}

// CreateJob saves the uploaded file and creates a pending job.
func (u *MediaUseCase) CreateJob(
	ctx context.Context,
	input UploadInput,
) (*domain.Job, error) {

	if input.File == nil || input.FileSize <= 0 {
		return nil, ErrNoFile
	}

	jobID := uuid.NewString()

	originalDir := filepath.Join(
		u.config.UploadDirectory,
		jobID,
	)

	if err := os.MkdirAll(originalDir, 0755); err != nil {
		return nil, fmt.Errorf(
			"create upload directory: %w",
			err,
		)
	}

	safeFileName := filepath.Base(input.FileName)

	originalPath := filepath.Join(
		originalDir,
		safeFileName,
	)

	dest, err := os.Create(originalPath)
	if err != nil {
		os.RemoveAll(originalDir)

		return nil, fmt.Errorf(
			"create file: %w",
			err,
		)
	}

	defer dest.Close()

	if _, err := io.Copy(dest, input.File); err != nil {
		os.RemoveAll(originalDir)

		return nil, fmt.Errorf(
			"save file: %w",
			err,
		)
	}

	now := time.Now()

	job := &domain.Job{
		ID:           jobID,
		OriginalName: safeFileName,
		OriginalPath: originalPath,
		OriginalSize: input.FileSize,
		Width:        input.Width,
		Quality:      input.Quality,
		Formats:      input.Formats,
		Status:       domain.StatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.repository.Create(ctx, job); err != nil {
		os.RemoveAll(originalDir)

		return nil, fmt.Errorf(
			"create job in repository: %w",
			err,
		)
	}

	return job, nil
}

// ProcessJob processes an existing pending job.
//
// This method implements worker.JobProcessor.
func (u *MediaUseCase) ProcessJob(
	ctx context.Context,
	jobID string,
) error {

	// Get the job.
	job, err := u.repository.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf(
			"get job %s: %w",
			jobID,
			err,
		)
	}

	// Make sure the job is still pending.
	if job.Status != domain.StatusPending {
		return fmt.Errorf(
			"job %s is not pending: %s",
			jobID,
			job.Status,
		)
	}

	// Mark job as processing.
	job.Status = domain.StatusProcessing
	job.UpdatedAt = time.Now()

	if err := u.repository.Update(ctx, job); err != nil {
		return fmt.Errorf(
			"mark job %s as processing: %w",
			jobID,
			err,
		)
	}

	// Output directory for this job (temporary — deleted after Cloudinary upload).
	outputDirectory := filepath.Join(
		u.config.OutputDirectory,
		jobID,
	)

	if err := os.MkdirAll(outputDirectory, 0755); err != nil {
		return u.failJob(
			ctx,
			job,
			fmt.Errorf(
				"create output directory: %w",
				err,
			),
		)
	}

	// Process the image locally first.
	localOutputs, _, err := u.processor.Process(
		ctx,
		processor.ProcessInput{
			JobID:           job.ID,
			OriginalPath:    job.OriginalPath,
			OutputDirectory: u.config.OutputDirectory,
			Width:           job.Width,
			Quality:         job.Quality,
			Formats:         job.Formats,
		},
	)

	if err != nil {
		os.RemoveAll(outputDirectory)
		return u.failJob(
			ctx,
			job,
			fmt.Errorf(
				"process job %s: %w",
				jobID,
				err,
			),
		)
	}

	// Upload each processed file to Cloudinary, then delete the local copy.
	cloudOutputs := make([]domain.Output, 0, len(localOutputs))

	for _, localOut := range localOutputs {
		publicID := fmt.Sprintf(
			"%s_optimized_%s",
			jobID,
			localOut.Format,
		)
		secureURL, fullPublicID, uploadErr := u.cloud.Upload(ctx, localOut.Path, publicID)
		if uploadErr != nil {
			// Non-fatal: log and continue with remaining formats.
			log.Printf("[ProcessJob] Cloudinary upload failed for %s (%s): %v", jobID, localOut.Format, uploadErr)
			continue
		}

		// Remove local file immediately after a successful upload.
		if err := os.Remove(localOut.Path); err != nil && !os.IsNotExist(err) {
			log.Printf("[ProcessJob] Failed to remove local file %s: %v", localOut.Path, err)
		}

		cloudOutputs = append(cloudOutputs, domain.Output{
			Format:   localOut.Format,
			Path:     secureURL,
			PublicID: fullPublicID,
			Size:     localOut.Size,
			Width:    localOut.Width,
			Height:   localOut.Height,
		})
	}

	// Clean up the (now empty) output directory and original upload.
	os.RemoveAll(outputDirectory)
	os.RemoveAll(filepath.Join(u.config.UploadDirectory, jobID))

	if len(cloudOutputs) == 0 {
		return u.failJob(ctx, job, fmt.Errorf("all Cloudinary uploads failed for job %s", jobID))
	}

	// Save processing results.
	job.Outputs = cloudOutputs
	job.Status = domain.StatusCompleted
	job.Error = ""
	job.UpdatedAt = time.Now()

	if err := u.repository.Update(ctx, job); err != nil {
		return fmt.Errorf(
			"complete job %s: %w",
			jobID,
			err,
		)
	}

	return nil
}

// failJob marks a job as failed and records the error.
func (u *MediaUseCase) failJob(
	ctx context.Context,
	job *domain.Job,
	err error,
) error {

	job.Status = domain.StatusFailed
	job.Error = err.Error()
	job.UpdatedAt = time.Now()

	if updateErr := u.repository.Update(ctx, job); updateErr != nil {
		return fmt.Errorf(
			"%w; update failed job: %v",
			err,
			updateErr,
		)
	}

	return err
}
