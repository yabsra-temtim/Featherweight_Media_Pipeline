package dto

import "featherweight/internal/domain"

// FromJob converts a domain Job into its wire representation.
func FromJob(job *domain.Job) JobResponse {
	return JobResponse{
		JobID:             job.ID,
		FileName:          job.OriginalName,
		OriginalSizeBytes: job.OriginalSize,
		Status:            string(job.Status),
		Outputs:           FromOutputs(job.Outputs),
		Error:             job.Error,
		Width:             job.Width,
		Quality:           job.Quality,
		Formats:           job.Formats,
		CreatedAt:         job.CreatedAt,
		UpdatedAt:         job.UpdatedAt,
	}
}

// FromOutput converts a domain Output into its wire representation.
func FromOutput(output domain.Output) OutputResponse {
	return OutputResponse{
		Format:    output.Format,
		URL:       output.Path,
		SizeBytes: output.Size,
		Width:     output.Width,
		Height:    output.Height,
	}
}

// FromOutputs converts a slice of domain Outputs into their wire representation.
// Always returns a non-nil slice so the JSON field serializes as `[]` rather
// than `null` when a job has no outputs yet.
func FromOutputs(outputs []domain.Output) []OutputResponse {
	result := make([]OutputResponse, 0, len(outputs))
	for _, output := range outputs {
		result = append(result, FromOutput(output))
	}
	return result
}
