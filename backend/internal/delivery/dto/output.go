package dto

// OutputResponse is the wire representation of a single processed output
// (one format/variant of an optimized image).
type OutputResponse struct {
	Format    string `json:"format"`
	URL       string `json:"url"`
	SizeBytes int64  `json:"size_bytes"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}
