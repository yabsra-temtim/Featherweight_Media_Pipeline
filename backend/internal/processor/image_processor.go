package processor

import (
	"context"
	"featherweight/internal/config"
	"featherweight/internal/domain"
)

// ProcessInput holds the data needed to process an image
type ProcessInput struct {
	JobID           string
	OriginalPath    string
	OutputDirectory string
	Width           int
	Quality         int
	Formats         []string
}

// ImageProcessor handles resizing and format conversion
type ImageProcessor struct {
	config config.Config
}

func NewImageProcessor(cfg config.Config) *ImageProcessor {
	return &ImageProcessor{
		config: cfg,
	}
}

// WebpEnabled tells the health route if we can process WebP
func (p *ImageProcessor) WebpEnabled() bool {
	return true
}

// AvifEnabled tells the health route if we can process AVIF
func (p *ImageProcessor) AvifEnabled() bool {
	return false
}

// Process does the actual image processing.
// For now, this is a skeleton so our Use Cases can compile!
func (p *ImageProcessor) Process(ctx context.Context, input ProcessInput) ([]domain.Output, []string, error) {
	// We will implement the actual image resizing logic later.
	// For now, let's just pretend it successfully created an output.

	outputs := []domain.Output{
		{
			Format: "jpeg",
			Path:   input.OutputDirectory + "/" + input.JobID + "/output.jpg",
			Size:   1024,
			Width:  input.Width,
		},
	}

	notes := []string{"Processing mocked for now"}

	return outputs, notes, nil
}
