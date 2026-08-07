package usecase

import (
	"errors"
	"io"

	"featherweight/internal/config"
	"featherweight/internal/domain"
	"featherweight/internal/processor"
)

var (
	ErrNoFile       = errors.New("an image file is required")
	ErrInvalidWidth = errors.New("width must be between 1 and the configured maximum")
)

// UploadInput represents the data we expect from the HTTP handler
type UploadInput struct {
	FileName string
	FileSize int64
	File     io.Reader
	Width    int
	Quality  int
	Formats  []string
}

// MediaUseCase orchestrates the upload and processing flow
type MediaUseCase struct {
	config     config.Config
	repository domain.JobRepository
	processor  *processor.ImageProcessor
}

func NewMediaUseCase(cfg config.Config, repo domain.JobRepository, proc *processor.ImageProcessor) *MediaUseCase {
	return &MediaUseCase{
		config:     cfg,
		repository: repo,
		processor:  proc,
	}
}
