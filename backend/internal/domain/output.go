package domain

type Output struct {
	Format string `json:"format"`
	path   string `json:"path"`
	size   int64  `json:"size"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}
