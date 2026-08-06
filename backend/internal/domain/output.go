package domain

// Output holds the result metadata for processed media.
type Output struct {
	Format string
	Path   string
	Size   int64
	Width  int
	Height int
}
