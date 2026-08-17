package dto

import "time"

// JobResponse is the wire representation of a media processing job,
// returned from both the upload endpoint and the job-status endpoint.
type JobResponse struct {
	JobID             string           `json:"job_id"`
	FileName          string           `json:"file_name"`
	OriginalSizeBytes int64            `json:"original_size_bytes"`
	Status            string           `json:"status"`
	Outputs           []OutputResponse `json:"outputs"`
	Error             string           `json:"error,omitempty"`
	Width             int              `json:"width"`
	Quality           int              `json:"quality"`
	Formats           []string         `json:"formats"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}
